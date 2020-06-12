use test;
create table user(
    id int primary key auto_increment,
    name varchar(100) default '',
    account varchar(100) default '',
    pwd varchar(100) default '',
    roles tinyint default 0,
    phone varchar(100) default '',
    email varchar(100) default '',
    company varchar(100) default '',
    status tinyint default 0,
    create_time datetime  default current_timestamp,
    update_time datetime  default current_timestamp on update current_timestamp
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
INSERT INTO `user`(`id`, `name`, `account`, `pwd`, `roles`, `phone`, `email`, `company`, `status`) VALUES (1, 'type', 'type@test.com', 'e10adc3949ba59abbe56e057f20f883e', 1, '', '', 'afg', 1);
