package model

// BuildDefaultGameDefinition
func BuildDefaultGameDefinition() GameDefinition {

	gd := GameDefinition{}
	gd.Duration = 1200000
	gd.MinParticipants = 2
	gd.MaxParticipants = 20
	gd.ArenaWidth = 2400
	gd.ArenaHeight = 1200
	gd.BulletSize = 16
	gd.LuchadorSize = 60
	gd.Fps = 30
	gd.BuletSpeed = 120

	gd.RadarAngle = 45
	gd.RadarRadius = 200
	gd.PunchAngle = 90
	gd.Life = 20
	gd.Energy = 30
	gd.PunchDamage = 2
	gd.PunchCoolDown = 2
	gd.MoveSpeed = 50
	gd.TurnSpeed = 180
	gd.TurnGunSpeed = 60
	gd.RespawnCooldown = 10
	gd.MaxFireCooldown = 10
	gd.MinFireDamage = 1
	gd.MaxFireDamage = 10
	gd.MinFireAmount = 1
	gd.MaxFireAmount = 10
	gd.RestoreEnergyperSecond = 3
	gd.RecycledLuchadorEnergyRestore = 6
	gd.IncreaseSpeedEnergyCost = 10
	gd.IncreaseSpeedPercentage = 20
	gd.FireEnergyCost = 2

	gd.RespawnX = 0
	gd.RespawnY = 0
	gd.RespawnAngle = 0
	gd.RespawnGunAngle = 0

	gd.GameComponents = make([]GameComponent, 0)
	gd.SceneComponents = make([]SceneComponent, 0)
	gd.Codes = make([]Code, 0)
	gd.LuchadorSuggestedCodes = make([]Code, 0)

	gd.MinLevel = 0
	gd.MaxLevel = 0
	gd.UnblockLevel = 0

	return gd

}
