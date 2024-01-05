Output post.gif
Set Width 1800
Set Height 1200
Set TypingSpeed 0.1
Set PlaybackSpeed 0.5

Hide
Type "kuroneko"
Enter
Sleep 1
Show

Sleep 1

# method
Enter
Down@.2
Enter

# URL
Tab@.2
Type 'https://httpbin.org/post'
Sleep 1

# content-type
Tab@.2
Enter
Down@.2
Enter

# body
Tab@.2 2
Enter
Sleep 2
Type 'i'
Type@10ms '{"foo":1, "bar": "baz"}'
Escape
Type ':wq'
Enter
Sleep 2

# Send
Tab@.2 1
Enter

Sleep 5
