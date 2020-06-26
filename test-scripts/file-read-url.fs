file::read "https://jsonplaceholder.typicode.com/posts" -> $source
log "downloaded in {source.took}"

loop $post in parse::json $source.data {
  if $post.title =~ /dolor/ {
    log "    {post.id}: {post.title}"
  }
}
