package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Bullet struct {
	Pos      rl.Vector2
	Velocity rl.Vector2
	LifeTime float32
}

type Particle struct {
	Pos        rl.Vector2
	Speed      rl.Vector2
	Color      rl.Color
	Radius     float32
	Active     bool
	ColorAlpha float32
}

type Enemy struct {
	Pos  rl.Vector2
	Dead bool
}

type Player struct {
	Health float32
	Pos    rl.Vector2
	Dead   bool
}

const (
	SCREEN_HEIGHT    = 500
	SCREEN_WIDTH     = 500
	FONT_SIZE        = 20
	TUTORIAL_TEXT    = "Press WASD to Move Around and LEFT click to shoot\n"
	PLAYER_DEAD_TEXT = "You Are Dead :) press F to restart"
	RELOAD_TEXT      = "Press R to Reload\n"
	BULLET_CAPACITY  = 30
	MAX_PARTICLES    = 30
	PARTICLE_SPEED   = 15
	MAX_TRAIL_POINTS = 20
	TRAIL_FADE_ALPHA = 20
)

var (
	lastEnemySpawn     = float64(0)
	enemySpawnInterval = float64(0.3)
	playerRadius       = float32(30)
	enemyRadius        = float32(30)
	bulletRadius       = float32(10)
	dx                 = float32(1000)
	dy                 = float32(1000)
	textAlpha          = float32(0)
	reloadTextAlpha    = float32(0)
	deadTextAlpha      = float32(0)
	showTutorial       = true
	bullets            = []Bullet{}
	player             = Player{
		Pos: rl.Vector2{
			X: float32(playerRadius + 10),
			Y: float32(playerRadius + 10),
		},
		Health: 100,
		Dead:   false,
	}

	enemies           []Enemy
	particles         [MAX_PARTICLES]Particle
	enemySpeed        = float32(0.02)
	playerColor       = rl.Color{230, 55, 55, 255}
	ballColor         = rl.Color{20, 190, 190, 255}
	trailPositions    [MAX_TRAIL_POINTS]rl.Vector2
	currentTrailIndex = 0
)

func manageTextAlpha(action string, dt float32) {
	if action == "fadeOut" {
		textAlpha -= 0.02
		if textAlpha <= 0 {
			textAlpha = 0
		}
		showTutorial = false
	}

	if action == "fadeIn" {
		textAlpha += dt
		if textAlpha >= 1 {
			textAlpha = 1
		}
	}
}

func manageDeadTextAlpha(dt float32) {
	if player.Dead {
		deadTextAlpha += dt
		if deadTextAlpha >= 1 {
			deadTextAlpha = 1
		}
	}

}

func shoot(camera rl.Camera2D) {
	mousePos := rl.GetMousePosition()
	mousePos = rl.GetScreenToWorld2D(mousePos, camera)
	playerPos := player.Pos
	bullet := Bullet{}
	direction := rl.Vector2Subtract(mousePos, playerPos)
	direction = rl.Vector2Normalize(direction)
	bullet.Velocity = rl.Vector2Scale(direction, 10)
	bullet.Pos = playerPos
	bullet.LifeTime = 3

	if len(bullets) < BULLET_CAPACITY {
		bullets = append(bullets, bullet)
	} else {
		// for showing reload text
		reloadTextAlpha = 1
	}
}

func spawnEnemies() {
	playerPos := player.Pos
	enemyPos := rl.Vector2{
		X: float32(rl.GetRandomValue(-800, 800)) + playerPos.X,
		Y: float32(rl.GetRandomValue(-800, 800)) + playerPos.Y,
	}
	enemies = append(enemies, Enemy{
		Pos:  enemyPos,
		Dead: false,
	})
}

func updateBulletPos(dt float32) {
	for i := 0; i < len(bullets)-1; i++ {
		bullets[i].Pos.X += bullets[i].Velocity.X
		bullets[i].Pos.Y += bullets[i].Velocity.Y

		bullets[i].LifeTime -= dt
	}
}

func updateEnemiesMovement() {
	for i := 0; i < len(enemies); i++ {
		enemies[i].Pos.X += (float32(player.Pos.X) - enemies[i].Pos.X) * enemySpeed
		enemies[i].Pos.Y += (float32(player.Pos.Y) - enemies[i].Pos.Y) * enemySpeed
	}
}

func initParticleBurst(origin rl.Vector2) {
	for i := 0; i < MAX_PARTICLES; i++ {
		particles[i].Pos = origin
		particles[i].Speed = rl.Vector2{
			X: rand.Float32() * PARTICLE_SPEED,
			Y: rand.Float32() * PARTICLE_SPEED,
		}
		particles[i].ColorAlpha = 1
		particles[i].Color = rl.Green
		particles[i].Radius = rand.Float32() * 20
		particles[i].Active = true

	}
}

func updateParticles() {
	for i := 0; i < MAX_PARTICLES; i++ {
		if particles[i].Active {
			particles[i].Pos.X += particles[i].Speed.X
			particles[i].Pos.Y += particles[i].Speed.Y
			particles[i].Radius -= 0.1
			particles[i].ColorAlpha -= 0.02
			particles[i].Color = rl.Fade(rl.Green, particles[i].ColorAlpha)
		}

		if particles[i].Radius <= 0.0 {
			particles[i].Active = false
		}
	}
}

func checkBulletEnemyCollision() {
	for i := 0; i < len(bullets); i++ {
		for j := 0; j < len(enemies); j++ {
			bulletPos := rl.Vector2{
				X: bullets[i].Pos.X,
				Y: bullets[i].Pos.Y,
			}
			enemyPos := rl.Vector2{
				X: enemies[j].Pos.X,
				Y: enemies[j].Pos.Y,
			}

			isCollide := rl.CheckCollisionCircles(bulletPos, bulletRadius, enemyPos, enemyRadius)

			if isCollide && !enemies[j].Dead && bullets[i].LifeTime >= 0 {

				if playerColor.G > 55 && playerColor.B > 55 {
					playerColor.G -= 20
					playerColor.B -= 20
				}

				if ballColor.R > 20 {
					ballColor.R += 16
				}

				if player.Health < 100 {
					player.Health += 10
				}

				initParticleBurst(enemyPos)
				enemies[j].Dead = true
				bullets[i].LifeTime = 0
			}
		}
	}
}

func checkEnemyPlayerCollision() {
	for i := 0; i < len(enemies); i++ {
		enemyPos := rl.Vector2{
			X: enemies[i].Pos.X,
			Y: enemies[i].Pos.Y,
		}

		isCollide := rl.CheckCollisionCircles(enemyPos, enemyRadius, player.Pos, playerRadius)

		if isCollide && player.Health > 0 && !enemies[i].Dead {
			playerColor.G += 20
			playerColor.B += 20
			ballColor.R += 16
			player.Health -= 10
			enemies[i].Dead = true
			initParticleBurst(enemyPos)

			if player.Health <= 0 {
				player.Dead = true
				initParticleBurst(player.Pos)
			}
		}

	}
}

func main() {

	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "circle game")
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetTargetFPS(60)

	camera := rl.Camera2D{
		Target: player.Pos,
		Offset: rl.Vector2{
			X: SCREEN_WIDTH / 2,
			Y: SCREEN_HEIGHT / 2,
		},
		Rotation: 0,
		Zoom:     1,
	}

	defer rl.CloseWindow()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()

		currentTime := rl.GetTime()

		if currentTime-lastEnemySpawn >= enemySpawnInterval && !player.Dead {
			spawnEnemies()
			lastEnemySpawn = currentTime
		}

		if showTutorial {
			manageTextAlpha("fadeIn", dt)
		}

		if rl.IsKeyDown(rl.KeyW) {
			manageTextAlpha("fadeOut", dt)
			player.Pos.Y -= float32(dy * dt)
		}
		if rl.IsKeyDown(rl.KeyS) {
			manageTextAlpha("fadeOut", dt)
			player.Pos.Y += float32(dy * dt)
		}

		if rl.IsKeyDown(rl.KeyD) {
			manageTextAlpha("fadeOut", dt)
			player.Pos.X += float32(dx * dt)
		}

		if rl.IsKeyDown(rl.KeyA) {
			manageTextAlpha("fadeOut", dt)
			player.Pos.X -= float32(dx * dt)
		}

		if rl.IsKeyDown(rl.KeyR) && len(bullets) == BULLET_CAPACITY {
			reloadTextAlpha = 0
			bullets = []Bullet{}
		}

		camera.Target = player.Pos

		//handle trail
		{
			trailPositions[currentTrailIndex] = rl.Vector2{X: player.Pos.X, Y: player.Pos.Y}
			currentTrailIndex = (currentTrailIndex + 1) % MAX_TRAIL_POINTS
		}

		// handle restart game
		if player.Dead && rl.IsKeyDown((rl.KeyF)) {
			player.Dead = false
			player.Health = 100
			playerColor = rl.Color{230, 55, 55, 255}
			ballColor = rl.Color{20, 190, 190, 255}
		}

		if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
			shoot(camera)
		}

		if !player.Dead {
			updateBulletPos(dt)
			updateEnemiesMovement()
			checkBulletEnemyCollision()
			checkEnemyPlayerCollision()
		}

		updateParticles()
		manageDeadTextAlpha(dt)

		if len(bullets) == BULLET_CAPACITY {
			reloadTextAlpha += dt
			if reloadTextAlpha >= 1 {
				reloadTextAlpha = 1
			}
		}

		rl.BeginDrawing()
		{
			rl.ClearBackground(rl.Black)

			rl.BeginMode2D(camera)
			{
				for i := 0; i < len(bullets)-1; i++ {
					if bullets[i].LifeTime >= 0 {
						rl.DrawCircle(int32(bullets[i].Pos.X), int32(bullets[i].Pos.Y), bulletRadius, rl.SkyBlue)
					}
				}

				for i := 0; i < len(enemies); i++ {
					if !enemies[i].Dead {
						rl.DrawCircle(int32(enemies[i].Pos.X), int32(enemies[i].Pos.Y), enemyRadius, ballColor)
					}
				}

				for i := 0; i < MAX_PARTICLES; i++ {
					if particles[i].Active {
						rl.DrawCircleV(particles[i].Pos, particles[i].Radius, particles[i].Color)
					}
				}

				if !player.Dead {
					rl.DrawCircle(int32(player.Pos.X), int32(player.Pos.Y), playerRadius, playerColor)

					for i := 0; i < MAX_TRAIL_POINTS; i++ {
						trailColor := playerColor
						trailColor.A = uint8(TRAIL_FADE_ALPHA * (MAX_TRAIL_POINTS - i) / MAX_TRAIL_POINTS)
						rl.DrawCircle(int32(trailPositions[i].X), int32(trailPositions[i].Y), playerRadius, trailColor)
					}
				}
			}

			rl.EndMode2D()

			rl.DrawText(RELOAD_TEXT, SCREEN_WIDTH/3, SCREEN_HEIGHT/2, FONT_SIZE, rl.Fade(rl.White, reloadTextAlpha))
			rl.DrawText(TUTORIAL_TEXT, SCREEN_WIDTH/5, SCREEN_HEIGHT/2, FONT_SIZE/3, rl.Fade(rl.White, textAlpha))
			if player.Dead {
				rl.DrawText(PLAYER_DEAD_TEXT, SCREEN_WIDTH/7, SCREEN_HEIGHT/2, FONT_SIZE, rl.Fade(rl.White, deadTextAlpha))
			}
		}
		rl.EndDrawing()
	}
}
