package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptAndParseCookie(t *testing.T) {
	key := "1gjrlcjQ8RyKANngp9607txr5fF5fhf1"
	input := "EQiu+j258LrdyiRYZxBiwCrqM+KLAWbVNJ0JZZpCfM9QdJUUWWj+sbBhdgCSfeuu/3o7MgSZVkm0m+Hk3YNxR6FKZxPN1jOG6N+0VGV7Mv5GpJePYZHRY3lTugrvuWzEleR7E0lXCdnbhkJjTsU+YeHTWdgqVuK34JcO5htwnzcgXxhFfY36Wj3GiM9RIbC6BdSVbKxHFuwbMatL99ffu3v+TOSR4IIi2ivewh1LUCIY/fOOPHuiwIs4XB+Qalf/B966p0tuEb1oG9kld8+XT4V/DgH8mIolBhqxnqVxTIXfCDyscZT/J3v7llf2vXUelCRVEV9M6JdXZxpe8fRT8Cv1+NgWP4bTKGfZMW3dFsos8IhmxhYibCrCP+mfbTPUgjJNIIEzfvxzsp9TpkIcxxVxgRx6GHyiUOvoLoNL6nbQ5cc2ZjyxXMOwi8XxYVOTaoEdL5oDQAzmSWfrUEDy4nRAzo5xUGKp/n5RMzViIT9G32fxCWoXBEk8kakNYPWFTaOGhDf00tyEnsiyCZTQkB5JwZGsOuxe++1Q/AC3RpY5+2fmS19S5rPvsyoKpEjcu2UytmMwR+C/nUKbWrrIEydaVJ5bwJKjmDraS3OPAH5EZg3UH1gVjrd3vlmLCGthITrh4/uUANrs3T4TQV0JaWsKWw0veb/8dFsRsc/Tp8r4FOSvAuL4tGarN/ObDcoafB/7+0dbLKzfPWlZi4QODmjK0bkjH3ptH/TBml5xWRNl62eg2lqYDU98niy1iV9xq65I0rv4ll2wSRrnlq0dSC1ERgxl3JC06Px52Hw+qDolNhoq7l29rfUuQVEkQ6xu7YM2EocWNWkeAPHUcOO8SSU37pBhHGFgNkzgzzM7LYHHoC1vtdftkYOerOGmwEFfQfN09JOXLDEVrDj1TmUXPKUZe/qozTeXN46by4eBzK4R+55Ri8T9t9tHPFsXnghIvue930A+UyHVVLurLgfskRA1IJBVLo/znf+IgYJuBwKZyB6drdmHKnHzt6w6MXxUK56NGKueTOY0Rv0QIIVSGvNaWWQYWxeB1nB3ac/g/W2BgxhHQDfFRS0QWAbGB8Rp6Hoxwu1FpcXF2GA8te4jVWAZNk7ijVAmtcnlMrIR0gBUMLO5n0tuvoS1cNtf2HlaReoXaGtlyDNmdvbiTgNAOUPqZ99iksCjEIaAT4axR6TZe2lbOax9LHMhbPCj1PWk9P8Z6iQDP0jWBbSOd7vfrPpBwkd9cxxIxkVIb4rwcInLmiJWiLLlKxigzg0H3X7hgVwNcbrpaLrikwJFqNMzK6GQ0o6MvQ30LqPW0EjUw+jir5cjGTHU2BnAm0pJf+Au/EBaUVguKinnBllXXo631+FLIIBqWdUsRN+H2oWSBtkukyN2mqjOlMCpIMjWBhZxvX3WhBvAjsJGKwk/WHM8rD3Rhqpant3E3aQnec62H0iP6mI8EixMYuO9QyaSDpWCps+CBOf4MCkhLmvGrFTpHUD/+OQMjxgx9msvctUJpB6pVlTSiyRgqagSU4LjQEsAFbHwxhzrEb2wanu+B1xCGXSnXZB9cQCkbgTVZ1IlYICL8XFHSHEE01OORD1eufghNulPbo5oOtGgj/eA3WK1vp5+M5vcfV/YexR9B7di5T7/ehXyNd8fCXSufPljo1xjo0c8CRxHbIvodf2MQI9Re6/KeQUlPhsFREt3aINxWpisciTjz8wA45SX153dgH+q0+ZpHd3lIDMCEIeULlWJmuzdz2Of1xIZvMbzl8tkpR+wZloX4F4WxUUXXnvBFX4Cc8uJsyvrTYD2n3LrZvb5WOL63jWnsWPBF8f6RLUMIDXnmGXSzFa9qk49u7Y+3hTd49wf8LjO6+FFbEbwkY7WX+zRKhO2gTA+AJHF/w5tzqQ1GPF0z3l7dOROoyP1CFJ41kq0Q1LaO7E8wpaxXU9NSQmMTzM89ag"

	user, _ := GetUser(input, key)
	assert.Equal(t, "Test Test", user.Name)
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, false, user.EmailVerified)
	assert.Equal(t, "Test", user.FirstName)
	assert.Equal(t, "Test", user.LastName)
	assert.Equal(t, "test@test.com", user.Email)
}
