# FastF1 Wrapper

A specialized adapter for the FastF1 library, bridging it to the Go backend.

# Good-to-Knows

## Track Layout Data

Imagine an F1 car has a "magic chalk" attached to its bottom. As the driver races around the track, the car draws a continuous line on the ground. Here is how we turn that line into the map you see:

### 1. What are these X and Y numbers?
Think of the race track as a giant piece of graph paper laid out over the real world. 
- **X** is how many steps the car moved **Left or Right**.
- **Y** is how many steps the car moved **Up or Down**.

The numbers you see (like `-1240.3`) represent the car's position in **decimeters** (about the length of a large smartphone). So if a car moves from `X: 0` to `X: 10`, it just moved 1 meter to the right.

### 2. Does FastF1 calculate the points?
Not exactly—the **cars** do! Every F1 car is full of GPS and movement sensors. Thousands of times per second, the car "pings" its location. FastF1 collects all those tiny pings and groups them into a "lap." 

We specifically ask for the **Fastest Lap** because that’s when the driver is on the "perfect line." If we used a slow lap, the map might look wiggly because the driver was avoiding traffic or pit-stopping!

### 3. Why do we rotate the track?
When the car records those X and Y points, they are based on the car's own internal compass, which might not match a map. Without rotating, the track might look sideways or upside down.

FastF1 gives us a "Rotation" value (like `44 degrees`). We use a bit of math to spin all those dots around a center point so that the track is oriented **North-Up**. It’s like turning a coloring book around until the picture is straight so you can color it properly.

### 4. What was that "Y-Flip" about?
This is a quirky computer thing. 
- In **Math Class**, bigger "Y" numbers go **UP**.
- On a **Computer Screen**, bigger "Y" numbers go **DOWN**.

If we don't "flip" the Y numbers, the track looks like it's being viewed in a mirror. We multiply the Y numbers by `-1` to make sure "Up" on the real track is "Up" on your screen.

### 5. How does the SVG draw the picture?
SVG stands for "Scalable Vector Graphics," but you can think of it as a digital **Connect-the-Dots** game.

Instead of sending a heavy image file (like a photo), our API sends a list of instructions:
1. "Put your pen down at the first X, Y dot."
2. "Draw a straight line to the next dot."
3. "Keep going until you've hit all 500 dots."
4. "Close the loop by drawing a line back to the start."

Because the computer is just following "drawing instructions" rather than looking at pixels, you can zoom in as much as you want and the track will **never get blurry**. It stays perfectly sharp, whether it's on a tiny phone or a giant TV!
