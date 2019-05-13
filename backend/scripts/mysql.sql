CREATE DATABASE `gocomic` 
    DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE gocomic;

CREATE TABLE `comics` (                                           
  `id` int(11) NOT NULL AUTO_INCREMENT,                           
  `url` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,         
  `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,       
  `source` varchar(32) COLLATE utf8mb4_general_ci NOT NULL,       
  `cover` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,       
  `summary` text COLLATE utf8mb4_general_ci,                      
  `latest` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,  
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,         
  PRIMARY KEY (`id`),                                             
  UNIQUE KEY `uix_comics_title_source` (`title`,`source`),        
  UNIQUE KEY `uix_comics_url` (`url`)                             
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `users` (                                            
  `id` int(11) NOT NULL AUTO_INCREMENT,                           
  `email` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,        
  `username` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,     
  `password` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,    
  `avatar` varchar(128) COLLATE utf8mb4_general_ci DEFAULT NULL,  
  `admin` tinyint(1) NOT NULL DEFAULT '0',                        
  `blocked` tinyint(1) NOT NULL DEFAULT '0',                      
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,         
  PRIMARY KEY (`id`),                                             
  UNIQUE KEY `uix_users_email` (`email`),                         
  UNIQUE KEY `uix_users_username` (`username`)                    
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

CREATE TABLE `subscribers` (                                                                                                                    
  `user_id` int(11) NOT NULL,                                                                                                                   
  `comic_id` int(11) NOT NULL,                                                                                                                  
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,                                                                                       
  PRIMARY KEY (`user_id`,`comic_id`),                                                                                                           
  KEY `idx_subscribers_user_id` (`user_id`),                                                                                                    
  KEY `idx_subscribers_comic_id` (`comic_id`),                                                                                                  
  CONSTRAINT `subscribers_comic_id_comics_id_foreign` FOREIGN KEY (`comic_id`) REFERENCES `comics` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `subscribers_user_id_users_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT     
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;                                                                             

CREATE TABLE `chapters` (                                                                                                                   
  `id` int(11) NOT NULL AUTO_INCREMENT,                                                                                                     
  `url` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,                                                                                   
  `number` int(11) NOT NULL,                                                                                                                
  `title` varchar(128) COLLATE utf8mb4_general_ci NOT NULL,                                                                                 
  `path` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,                                                                                      
  `cached` tinyint(1) NOT NULL DEFAULT '0',                        
  `comic_id` int(11) NOT NULL,                                                                                                              
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,                                                                                   
  PRIMARY KEY (`id`),                                                                                                                       
  UNIQUE KEY `uix_chapters_url` (`url`),                                                                                                    
  UNIQUE KEY `uix_chapters_number_comic_id` (`number`,`comic_id`),                                                                          
  KEY `chapters_comic_id_comics_id_foreign` (`comic_id`),                                                                                   
  CONSTRAINT `chapters_comic_id_comics_id_foreign` FOREIGN KEY (`comic_id`) REFERENCES `comics` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;                                                                         

CREATE TABLE `pictures` (                                                                                                                           
  `id` bigint(20) NOT NULL AUTO_INCREMENT,                                                                                                          
  `number` int(11) NOT NULL,                                                                                                                
  `filename` varchar(256) COLLATE utf8mb4_general_ci DEFAULT NULL,                                                                                      
  `chapter_id` int(11) NOT NULL,                                                                                                                    
  `created` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,                                                                                           
  PRIMARY KEY (`id`),                                                                                                                               
  KEY `idx_pictures_chapter_id` (`chapter_id`),                                                                                                     
  CONSTRAINT `pictures_chapter_id_chapters_id_foreign` FOREIGN KEY (`chapter_id`) REFERENCES `chapters` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;                                                                                  
