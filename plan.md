The goal here is to provide a simple client-side frontend which calls the go rpc server.

1. The API surfaces threadded discussion.
2. Patchwork-style "you have 4 new things. click to update"
3. Initial page load fetches 10 posts.

Open questions:
- How do FE's authenticate with the backend such that they can get read access?
- What is the correct API to get a post list?

For later:
- subscribe to websocket w/ a filter so we don't have to process all data.

--

1. figure out how to get author name (get post by type for 'about')
1. figure out threadding
1. Get most recent posts (by claimed date), not the least recent posts
