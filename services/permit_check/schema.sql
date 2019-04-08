CREATE TABLE IF NOT EXISTS VEHICLE (
    vehicleid VARCHAR(6),
    licensenumber VARCHAR(10),
    statecode VARCHAR(2),
    PRIMARY KEY(vehicleid)
);

CREATE TABLE IF NOT EXISTS PERMIT (
    ownername VARCHAR(30),
    permitid VARCHAR(8),
    PRIMARY KEY(permitid)
);

CREATE TABLE IF NOT EXISTS VEH_PER (
    vid VARCHAR(6),
    pid VARCHAR(8),
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_p(pid) REFERENCES PERMIT(permitid),
    PRIMARY KEY(vid,pid)
);

CREATE TABLE IF NOT EXISTS LOT (
    lotid VARCHAR(36),
    name VARCHAR(64),
    occupancy INT,
    PRIMARY KEY(lotid)
);

CREATE TABLE IF NOT EXISTS LOT_PERMIT (
    pid VARCHAR(8),
    lid VARCHAR(36),
    FOREIGN KEY fk_p(pid) REFERENCES PERMIT(permitid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY (pid,lid)
);

CREATE TABLE IF NOT EXISTS ENTRY_RECORD (
    vid VARCHAR(6),
    lid VARCHAR(36),
    timestamp DATE NOT NULL,
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY(vid,lid,timestamp)
);

CREATE TABLE IF NOT EXISTS EXIT_RECORD (
    vid VARCHAR(6),
    lid VARCHAR(36),
    timestamp DATE NOT NULL,
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY(vid,lid,timestamp)
);

CREATE TABLE IF NOT EXISTS OCCUPANCY (
    lid VARCHAR(36),
    vid VARCHAR(6),
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY(lid,vid)
);

DELIMITER $$
CREATE TRIGGER enter_lot 
    BEFORE INSERT ON ENTRY_RECORD
    FOR EACH ROW
    BEGIN   
        INSERT INTO OCCUPANCY (lid,vid) VALUES (NEW.lid, NEW.vid);
    END; 
    $$
DELIMITER ;

DELIMITER $$
CREATE TRIGGER exit_lot 
    BEFORE INSERT ON EXIT_RECORD 
    FOR EACH ROW
    BEGIN
        DELETE FROM OCCUPANCY WHERE lid=NEW.lid AND vid = NEW.vid;
    END;
    $$
DELIMITER ;