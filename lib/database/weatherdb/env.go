package weatherdb

import "time"

type SensorType struct {
	Id          int    `json:"id" xml:"id,attr"`
	Type        string `json:"type" xml:"type"`
	Unit        string `json:"unit" xml:"unit"`
	Description string `json:"description" xml:"description"`
}

type Sensor struct {
	Id          int    `json:"id" xml:"id,attr"`
	Name        string `json:"name" xml:"name"`
	Type        string `json:"type" xml:"type"`
	TypeId      string `json:"type_id" xml:"type_id,attr"`
	Unit        string `json:"unit" xml:"unit"`
	Description string `json:"description" xml:"description"`
}

type SensorData struct {
	SensorId  int        `json:"sensor_id" xml:"sensor_id,attr"`
	Timestamp time.Time  `json:"timestamp" xml:"timestamp"`
	Name      string     `json:"name" xml:"name"`
	Type      string     `json:"type" xml:"type"`
	Unit      string     `json:"unit" xml:"unit"`
	Value     int64      `json:"value" xml:"value"`
}

type SensorValue struct {
	Timestamp time.Time `json:"timestamp" xml:"timestamp"`
	Value     int64     `json:"value" xml:"value"`
}

type Window struct {
	Image            string       `json:"image" xml:"image"`
	Title            string       `json:"title" xml:"title"`
	Timestamp        time.Time    `json:"timestamp" xml:"timestamp"`
	WindowStateImage string       `json:"window_state_image" xml:"window_state_image"`
	Open             bool         `json:"open" xml:"open"`
}

type Door struct {
	Image string `json:"image" xml:"image"`
	Title string `json:"title" xml:"title"`
}

type Subscription struct {
	SensorId      int `json:"sensor_id" xml:"sensor_id,attr"`
	UserId      int `json:"user_id" xml:"user_id,attr"`
}

type User struct {
	Id            int            `json:"id" xml:"id,attr"`
	Password      string         `json:"password" xml:"password"`
	Name          string         `json:"name" xml:"name"`
	Email         string         `json:"email" xml:"email"`
	Role          string         `json:"role" xml:"role"`
	Active        bool           `json:"active" xml:"active"`
	Subscriptions []Subscription `json:"subscriptions" xml:"subscriptions"`
}
