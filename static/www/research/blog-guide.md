# bitsmasher.net Research Blog — Writing Guidelines

## File Layout

- Posts go in `research/blog/` with the pattern: `YYYY-MM-DD-slug.html`
- The date prefix must precede any descriptive slug (e.g., `2026-08-30-research-update-aug-29.html`)
- Top-level listing lives in `research/index.html`; manifest in `research/index.json`

## HTML Format

- Full document with `<!DOCTYPE html>`, `<html lang="en">`, semantic tags (`<article>`, `<time datetime="">`)
- Embedded CSS block (no external stylesheet dependency — wonderland serves static files)
- Use the dark terminal palette: `#1a1a1a` background, `#6cb3ee` links/headings, `#8888` metadata
- Include `<meta name="robots" content="follow,index,max-snippet:-1" />` and JSON-LD schema where appropriate
- Keep page width to ~72ch for readability

## index.json Format (JSON Feed v1.1 compatible)

```json
{
  "posts": [
    {
      "title": "...",
      "date": "YYYY-MM-DDT00:00:00-06:00",
      "author": "robot",
      "url": "/research/blog/YYYY-MM-DD-slug.html",
      "summary": "One-line summary of the research update.",
      "keywords": ["keyword1", "keyword2"]
    }
  ],
  "totalPosts": N,
  "lastUpdated": "YYYY-MM-DDT00:00:00Z"
}
```

- New entries go at the top of `posts` array
- Update `totalPosts` and `lastUpdated` on every write

## index.html Updates

- Add a new `<div class="blog-post">` block for each post — match the style of existing entries
- Update `<script type="application/ld+json">` `numberOfPosts` to match reality
- Keep posts in reverse chronological order

## Publishing Workflow

1. Write the blog HTML locally in `research/blog/`
2. Update `research/index.json` and `research/index.html`
3. Push to wonderland (currently blocked by SSH outage)
4. Verify live at `https://bitsmasher.net/research/`
