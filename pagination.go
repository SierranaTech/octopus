package octopus

import (
	"context"
	"iter"
	"net/url"
)

const maxPages = 100

// paginate follows paginated responses and collects all results.
// Returns ErrPageLimitExceeded if the page cap is reached.
func paginate[T any](ctx context.Context, c *Client, path string, query url.Values, auth bool) ([]T, error) {
	var all []T

	var page Page[T]
	if err := c.get(ctx, path, query, auth, &page); err != nil {
		return nil, err
	}
	all = append(all, page.Results...)

	for i := 1; page.Next != nil && *page.Next != ""; i++ {
		if i >= maxPages {
			return all, ErrPageLimitExceeded
		}
		nextURL := *page.Next
		page = Page[T]{}
		if err := c.getURL(ctx, nextURL, auth, &page); err != nil {
			return nil, err
		}
		all = append(all, page.Results...)
	}

	return all, nil
}

// pages returns an iterator that yields one page of results at a time.
// This allows streaming processing without loading all pages into memory.
func pages[T any](ctx context.Context, c *Client, path string, query url.Values, auth bool) iter.Seq2[[]T, error] {
	return func(yield func([]T, error) bool) {
		var page Page[T]
		if err := c.get(ctx, path, query, auth, &page); err != nil {
			yield(nil, err)
			return
		}
		if !yield(page.Results, nil) {
			return
		}

		for i := 1; page.Next != nil && *page.Next != ""; i++ {
			if i >= maxPages {
				yield(nil, ErrPageLimitExceeded)
				return
			}
			nextURL := *page.Next
			page = Page[T]{}
			if err := c.getURL(ctx, nextURL, auth, &page); err != nil {
				yield(nil, err)
				return
			}
			if !yield(page.Results, nil) {
				return
			}
		}
	}
}
