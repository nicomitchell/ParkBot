CREATE TABLE VEHICLE (
    vehicleid VARCHAR(6),
    licensenumber VARCHAR(10),
    statecode VARCHAR(2),
    PRIMARY KEY(vehicleid)
);

CREATE TABLE PERMIT (
    ownername VARCHAR(30),
    permitid VARCHAR(8),
    PRIMARY KEY(permitid)
);

CREATE TABLE VEH_PER (
    vid VARCHAR(6),
    pid VARCHAR(8),
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_p(pid) REFERENCES PERMIT(permitid),
    PRIMARY KEY(vid,pid)
);

CREATE TABLE LOT (
    lotid VARCHAR(36),
    name VARCHAR(64),
    occupancy INT,
    PRIMARY KEY(lotid)
);

CREATE TABLE LOT_PERMIT (
    pid VARCHAR(8),
    lid VARCHAR(36)
    FOREIGN KEY fk_p(pid) REFERENCES PERMIT(permitid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY (pid,lid)
);

CREATE TABLE ENTRY_RECORD (
    vid VARCHAR(6),
    lid VARCHAR(36),
    timestamp DATE NOT NULL,
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY rec_id(vid,lid,timestamp)
);

CREATE TABLE EXIT_RECORD (
    vid VARCHAR(6),
    lid VARCHAR(36),
    timestamp DATE NOT NULL,
    FOREIGN KEY fk_v(vid) REFERENCES VEHICLE(vehicleid),
    FOREIGN KEY fk_l(lid) REFERENCES LOT(lotid),
    PRIMARY KEY rec_id(vid,lid,timestamp)
);
