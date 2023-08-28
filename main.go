package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type GameState int

const (
	Playing GameState = iota
	GameOver
)

type Player struct {
	rl.RectangleInt32
	direction direction
	Tail      []Tail
}

type direction struct {
	up, down, left, right bool
}

type Food struct {
	rl.RectangleInt32
}
type Tail struct {
	rl.RectangleInt32
}

var Score int = 0
var widthInt, heightInt int32 = 0, 0

func main() {
	rl.InitWindow(720, 720, "Snake test game")
	rl.SetTargetFPS(60)

	widthInt, heightInt = int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())

	for !rl.WindowShouldClose() {

		SnakeGame()
		replay := false

		for {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			rl.DrawText("Game Over", widthInt/2-50, heightInt/2, 20, rl.White)
			rl.DrawText("Score: "+fmt.Sprint(Score), widthInt/2-50, heightInt/2+20, 20, rl.White)
			rl.DrawText("Press Enter to replay", widthInt/2-50, heightInt/2+40, 20, rl.White)
			rl.DrawText("Press Esc to quit", widthInt/2-50, heightInt/2+60, 20, rl.White)

			if rl.IsKeyPressed(rl.KeyEnter) {
				replay = true
				Score = 0
				break
			} else if rl.IsKeyPressed(rl.KeyEscape) {
				break
			}

			rl.EndDrawing()
		}

		if !replay {
			break
		}
	}

	rl.CloseWindow()
}

func SnakeGame() {

	var player = Player{
		RectangleInt32: rl.RectangleInt32{
			X:      widthInt / 2,
			Y:      heightInt / 2,
			Width:  10,
			Height: 10,
		},
		direction: direction{false, false, false, true},
		Tail:      nil,
	}
	var food = Food{
		RectangleInt32: rl.RectangleInt32{
			X:      400,
			Y:      200,
			Width:  10,
			Height: 10,
		},
	}

	accumulator := float32(0.0)

	movementInterval := float32(0.1)

	var i, j int32
	gameState := Playing

	for gameState == Playing {
		accumulator += rl.GetFrameTime()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.DrawRectangleLines(0, 0, widthInt, heightInt, rl.Black)

		for i = 0; i < widthInt; i += 10 {
			rl.DrawLine(i, 0, i, heightInt, rl.Black)
		}

		for j = 0; j < heightInt; j += 10 {
			rl.DrawLine(0, j, widthInt, j, rl.Black)
		}

		rl.DrawRectangle(food.X, food.Y, food.Width, food.Height, rl.Green)
		rl.DrawRectangle(player.X, player.Y, player.Width, player.Height, rl.Pink)

		DrawTail(&player)

		if (rl.IsKeyDown(rl.KeyW) || rl.IsKeyDown(rl.KeyUp)) && !player.direction.down {
			player.direction.up = true
			player.direction.down = false
			player.direction.left = false
			player.direction.right = false
		}
		if (rl.IsKeyDown(rl.KeyS) || rl.IsKeyDown(rl.KeyDown)) && !player.direction.up {
			player.direction.up = false
			player.direction.down = true
			player.direction.left = false
			player.direction.right = false
		}
		if (rl.IsKeyDown(rl.KeyA) || rl.IsKeyDown(rl.KeyLeft)) && !player.direction.right {
			player.direction.up = false
			player.direction.down = false
			player.direction.left = true
			player.direction.right = false
		}
		if (rl.IsKeyDown(rl.KeyD) || rl.IsKeyDown(rl.KeyRight)) && !player.direction.left {
			player.direction.up = false
			player.direction.down = false
			player.direction.left = false
			player.direction.right = true
		}

		if accumulator >= movementInterval {

			if TailCollisionRecInt32(player.Tail, player.RectangleInt32) ||
				player.RectangleInt32.X < 0 || player.RectangleInt32.X > widthInt ||
				player.RectangleInt32.Y < 0 || player.RectangleInt32.Y > heightInt {
				gameState = GameOver
			}

			if CollisionRecInt32(player.RectangleInt32, food.RectangleInt32) {
				Score++
				movementInterval *= 0.9
				food.X = rl.GetRandomValue(1, widthInt-1) / 10 * 10
				food.Y = rl.GetRandomValue(1, heightInt-1) / 10 * 10
				AddTail(&player)
			}

			UpdateTail(&player)

			if player.direction.up {
				player.Y -= 10
			}
			if player.direction.down {
				player.Y += 10
			}
			if player.direction.left {
				player.X -= 10
			}
			if player.direction.right {
				player.X += 10
			}

			accumulator = 0.0
		}

		if player.X < 0 || player.X > widthInt || player.Y < 0 || player.Y > heightInt {
			gameState = GameOver
		}
		rl.EndDrawing()
	}
}

func TailCollisionRecInt32(tail []Tail, r2 rl.RectangleInt32) bool {
	for _, r1 := range tail {
		if r1.X < r2.X+r2.Width &&
			r1.X+r1.Width > r2.X &&
			r1.Y < r2.Y+r2.Height &&
			r1.Y+r1.Height > r2.Y {
			return true
		}
	}
	return false
}

func CollisionRecInt32(r1, r2 rl.RectangleInt32) bool {
	return r1.X < r2.X+r2.Width &&
		r1.X+r1.Width > r2.X &&
		r1.Y < r2.Y+r2.Height &&
		r1.Y+r1.Height > r2.Y
}

func AddTail(player *Player) {
	var tail = Tail{
		RectangleInt32: rl.RectangleInt32{
			X:      player.X,
			Y:      player.Y,
			Width:  10,
			Height: 10,
		},
	}
	player.Tail = append(player.Tail, tail)
}

func UpdateTail(player *Player) {
	if len(player.Tail) == 0 {
		return
	}
	for i := len(player.Tail) - 1; i > 0; i-- {
		player.Tail[i].X = player.Tail[i-1].X
		player.Tail[i].Y = player.Tail[i-1].Y
	}
	player.Tail[0].X = player.X
	player.Tail[0].Y = player.Y

}

func DrawTail(player *Player) {
	for _, tail := range player.Tail {
		rl.DrawRectangle(tail.X, tail.Y, tail.Width, tail.Height, rl.Color{255, 120, 175, 255})
	}
}
