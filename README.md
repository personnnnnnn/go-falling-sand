A little falling sand game (framework?) I am making in go.

You might be wondering, _why go?_
Well, 3 reasons:

- **It's fast**: Speed really matters when you are dealing with hunderds, if not thousands of cells that need to be updated and displayed to the screen every frame.
- **Goroutines**: For me, dealing with multithreading in the past with diffrent languages has not been a pleasant experience, to say the least.
- **It has a garbage collector**: I'm lazy. What can I say?

For graphics, I am using the amazing [ebiten](https://github.com/hajimehoshi/ebiten) library.

For now, this is it.
