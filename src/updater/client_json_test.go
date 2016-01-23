package updater

const (
	recLoadAllFailer = `
{
	"msg": "I am error",
	"result": "failed"
}
`

	recLoadAllSuccess = `
{
	"msg": "success",
	"result": "success",
	"response": {
		"recs": {
			"objs": [
				{
					"rec_id": "1",
					"content": "127.0.0.1",
					"type": "A",
					"name": "test"
				},
				{
					"rec_id": "2",
					"content": "::1",
					"type": "AAAA",
					"name": "test"
				}
			]
		}
	}
}
`

	recLoadAllSuccessIpV4 = `
{
	"msg": "success",
	"result": "success",
	"response": {
		"recs": {
			"objs": [
				{
					"rec_id": "1",
					"content": "127.0.0.1",
					"type": "A",
					"name": "test"
				}
			]
		}
	}
}
`

	recLoadAllSuccessIpV6 = `
{
	"msg": "success",
	"result": "success",
	"response": {
		"recs": {
			"objs": [
				{
					"rec_id": "2",
					"content": "::1",
					"type": "AAAA",
					"name": "test"
				}
			]
		}
	}
}
`

	editResFailer = `
{
	"result": "failed"
}
`

	editResSuccess = `
{
	"result": "success"
}
`

	shoddyJson = "I am shoddy json"
)
