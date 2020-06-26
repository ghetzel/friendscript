file::read "https://jsonplaceholder.typicode.com/posts" -> $source
parse::json $source.data -> $json

log $json