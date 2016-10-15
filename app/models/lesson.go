package models

import (
"time"
//"github.com/jinzhu/gorm"
)

/*
CREATE TABLE `lessons` (
    `id` int unsigned AUTO_INCREMENT,
    `created_at` timestamp NULL,
    `updated_at` timestamp NULL,
    `deleted_at` timestamp NULL,

    `post_id` int unsigned,
    `status` int,
    `start_time` timestamp NULL,
    `timeout` int unsigned , 
    `lng` double,
    `lat` double ,
PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8;
CREATE INDEX idx_lessons_deleted_at ON `lessons`(deleted_at);
*/
type Lesson struct{
    BaseModel
 
    TopicID      uint       `gorm:"not null"`
    Status      int
    StartTime   time.Time
    Timeout     uint
    Lng         float64
    Lat         float64
}