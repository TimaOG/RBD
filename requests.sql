CREATE TABLE users (
	pkId int AUTO_INCREMENT,
    Password text,
    EMail varchar(255),
    PhoneNumber varchar(255),
    IsAdmin tinyint(1),
    PRIMARY KEY (pkId)
);
CREATE TABLE ordersType (
    pkId int,
    typeName varchar(255),
    price int,
    discribtion text,
    PRIMARY KEY (pkId)
);
CREATE TABLE orders (
	pkId int AUTO_INCREMENT,
    fkCodeClient int,
    fkOrderType int,
    address text,
    usefulTime datetime,
    discribtion text,
    requestStatus varchar(255),
    startTime datetime,
    closeTime datetime,
    FOREIGN KEY (fkCodeClient) REFERENCES users (pkId),
    FOREIGN KEY (fkOrderType) REFERENCES ordersType (pkId),
    PRIMARY KEY (pkId)
);

