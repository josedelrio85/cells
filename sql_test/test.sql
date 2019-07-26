CREATE TABLE leadnew (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  legacy_id bigint(20) DEFAULT NULL,
  lea_ts timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  lea_smartcenter_id varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  sou_id bigint(20) DEFAULT NULL,
  leatype_id bigint(20) DEFAULT NULL,
  passport_id varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  passport_id_grp varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  utm_source varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  sub_source varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  lea_phone varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  lea_mail varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  lea_name varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  lea_url varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  lea_ip varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  is_smart_center tinyint(1) DEFAULT NULL,
  lea_dni varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  gclid varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  domain varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  observations text COLLATE utf8_spanish_ci,
  PRIMARY KEY (id),
  KEY idx_leadnew_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;


CREATE TABLE creditea (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  lea_id int(10) unsigned DEFAULT NULL,
  requestedamount varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  contracttype varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  netincome varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  outofschedule varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  asnef tinyint(1) DEFAULT '0',
  alreadyclient tinyint(1) DEFAULT '0',
  PRIMARY KEY (id),
  KEY idx_creditea_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;

CREATE TABLE rcableexp (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  lea_id int(10) unsigned DEFAULT NULL,
  location varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  answer varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  respvalues varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  coverture varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_rcableexp_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;


CREATE TABLE microsoft (
  id int(10) unsigned NOT NULL AUTO_INCREMENT,
  created_at timestamp NULL DEFAULT NULL,
  updated_at timestamp NULL DEFAULT NULL,
  deleted_at timestamp NULL DEFAULT NULL,
  lea_id int(10) unsigned DEFAULT NULL,
  tipoordenador varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  sector varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  presupuesto varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  rendimiento varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  movilidad varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  office365 varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  observaciones varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  producttype varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  productname varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  product_id varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  originalprice varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  price varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  brand varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  discountpercentage varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  discountcode varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  processortype varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  diskcapacity varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  graphics varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  wirelessinterface varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  devicesaverageage varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  devicesoperatingsystem varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  deviceshangfrequency varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  devicesnumber varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  deviceslastyearrepairs varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  devicesstartuptime varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  pageindex tinyint(1) DEFAULT NULL,
  oldsouid bigint(20) DEFAULT NULL,
  tipouso varchar(255) COLLATE utf8_spanish_ci DEFAULT NULL,
  PRIMARY KEY (id),
  KEY idx_microsoft_deleted_at (deleted_at)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci;



CREATE TABLE sources (
  sou_id int(3) NOT NULL AUTO_INCREMENT,
  sou_description varchar(50) CHARACTER SET utf8 COLLATE utf8_spanish_ci DEFAULT NULL,
  sou_idcrm int(3) DEFAULT NULL,
  PRIMARY KEY (sou_id)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci ROW_FORMAT=COMPACT;



CREATE TABLE leadtypes (
  leatype_id int(3) NOT NULL AUTO_INCREMENT,
  leatype_description varchar(50) CHARACTER SET utf8 COLLATE utf8_spanish_ci DEFAULT NULL,
  leatype_idcrm int(3) DEFAULT NULL,
  PRIMARY KEY (leatype_id)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COLLATE=utf8_spanish_ci ROW_FORMAT=COMPACT;




insert into sources (sou_id, sou_description, sou_idcrm) values (1,'CREDITEA ABANDONADOS',2);
insert into sources (sou_id, sou_description, sou_idcrm) values (2,' CREDITEA STAND',3);
insert into sources (sou_id, sou_description, sou_idcrm) values (3,'EVO BANCO',4);
insert into sources (sou_id, sou_description, sou_idcrm) values (4,'CREDITEA TIMEOUT',5);
insert into sources (sou_id, sou_description, sou_idcrm) values (5,'R CABLE',6);
insert into sources (sou_id, sou_description, sou_idcrm) values (6,'BYSIDECAR',7);
insert into sources (sou_id, sou_description, sou_idcrm) values (7,'HERCULES',8);
insert into sources (sou_id, sou_description, sou_idcrm) values (8,'SEGURO PARA MOVIL',11);
insert into sources (sou_id, sou_description, sou_idcrm) values (9,'CREDITEA END TO END',13);
insert into sources (sou_id, sou_description, sou_idcrm) values (10,'CREDITEA FB',14);
insert into sources (sou_id, sou_description, sou_idcrm) values (11,'CREDITEA RASTREATOR',15);
insert into sources (sou_id, sou_description, sou_idcrm) values (12,'EUSKALTEL',16);
insert into sources (sou_id, sou_description, sou_idcrm) values (13,'ADESLAS',19);
insert into sources (sou_id, sou_description, sou_idcrm) values (14,'R CABLE EMPRESAS',20);
insert into sources (sou_id, sou_description, sou_idcrm) values (15,'PRUEBA BySidecar',23);
insert into sources (sou_id, sou_description, sou_idcrm) values (16,'EVO BANCO FIRMADOS NO FORMALIZ',24);
insert into sources (sou_id, sou_description, sou_idcrm) values (17,'YOIGO NEGOCIOS DERIVACION YOIG',25);
insert into sources (sou_id, sou_description, sou_idcrm) values (18,'YOIGO NEGOCIOS SEO',26);
insert into sources (sou_id, sou_description, sou_idcrm) values (19,'YOIGO NEGOCIOS SEM',27);
insert into sources (sou_id, sou_description, sou_idcrm) values (20,'YOIGO NEGOCIOS EMAILING',28);
insert into sources (sou_id, sou_description, sou_idcrm) values (21,'CREDITEA C2C ATT. CLIENTE',29);
insert into sources (sou_id, sou_description, sou_idcrm) values (22,'CREDITEA C2C NOT STARTED',30);
insert into sources (sou_id, sou_description, sou_idcrm) values (23,'CREDITEA PAGO RECURRENTE',31);
insert into sources (sou_id, sou_description, sou_idcrm) values (24,'SANAL',32);
insert into sources (sou_id, sou_description, sou_idcrm) values (25,'MICROSOFT R CABLE',33);
insert into sources (sou_id, sou_description, sou_idcrm) values (26,'MICROSOFT RECOMENDADOR RRSS',34);
insert into sources (sou_id, sou_description, sou_idcrm) values (27,'MICROSOFT RECOMENDADOR  SEO',35);
insert into sources (sou_id, sou_description, sou_idcrm) values (28,'MICROSOFT RECOMENDADOR GOOGLE',36);
insert into sources (sou_id, sou_description, sou_idcrm) values (29,'MICROSOFT RECOMENDADOR EMAILIN',37);
insert into sources (sou_id, sou_description, sou_idcrm) values (30,'MICROSOFT RECOMENDADOR PROGRAM',38);
insert into sources (sou_id, sou_description, sou_idcrm) values (31,'MICROSOFT OFERTAS GOOGLE',39);
insert into sources (sou_id, sou_description, sou_idcrm) values (32,'MICROSOFT OFERTAS EMAILING',40);
insert into sources (sou_id, sou_description, sou_idcrm) values (33,'MICROSOFT OFERTAS  SEO',41);
insert into sources (sou_id, sou_description, sou_idcrm) values (34,'MICROSOFT OFERTAS PROGRAMÁTICA',42);
insert into sources (sou_id, sou_description, sou_idcrm) values (35,'MICROSOFT OFERTAS RRSS',43);
insert into sources (sou_id, sou_description, sou_idcrm) values (36,'MICROSOFT FICHA PRODUCTO  SEO',44);
insert into sources (sou_id, sou_description, sou_idcrm) values (37,'MICROSOFT FICHA PRODUCTO GOOGL',45);
insert into sources (sou_id, sou_description, sou_idcrm) values (38,'MICROSOFT FICHA PRODUCTO EMAIL',46);
insert into sources (sou_id, sou_description, sou_idcrm) values (39,'MICROSOFT FICHA PRODUCTO PROGR',47);
insert into sources (sou_id, sou_description, sou_idcrm) values (40,'MICROSOFT FICHA PRODUCTO RRSS',48);
insert into sources (sou_id, sou_description, sou_idcrm) values (41,'MICROSOFT PERDIDAS',49);
insert into sources (sou_id, sou_description, sou_idcrm) values (42,'MICROSOFT  PERDIDAS RECOMENDAD',51);
insert into sources (sou_id, sou_description, sou_idcrm) values (43,'MICROSOFT  PERDIDAS OFERTAS',52);
insert into sources (sou_id, sou_description, sou_idcrm) values (44,'MICROSOFT  PERDIDAS FICHA PROD',53);
insert into sources (sou_id, sou_description, sou_idcrm) values (45,'CREDITEA C2C DOMINGO',54);
insert into sources (sou_id, sou_description, sou_idcrm) values (46,'MICROSOFT HAZELCAMBIO',55);
insert into sources (sou_id, sou_description, sou_idcrm) values (47,'MICROSOFT HAZELCAMBIO PERDIDAS',56);
insert into sources (sou_id, sou_description, sou_idcrm) values (48,'MICROSOFT CALCULADORA',57);
insert into sources (sou_id, sou_description, sou_idcrm) values (49,'MICROSOFT RECOMENDADOR',58);
insert into sources (sou_id, sou_description, sou_idcrm) values (50,'MICROSOFT OFERTAS',59);
insert into sources (sou_id, sou_description, sou_idcrm) values (51,'MICROSOFT FICHA PRODUCTO',60);
insert into sources (sou_id, sou_description, sou_idcrm) values (52,'MICROSOFT GLOBAL',61);
insert into sources (sou_id, sou_description, sou_idcrm) values (53,'CREDITEA BO NOT STARTED',62);
insert into sources (sou_id, sou_description, sou_idcrm) values (54,'R CABLE EXPANSION END TO END',63);
insert into sources (sou_id, sou_description, sou_idcrm) values (55,'R CABLE EXPANSION ENTRANTE',64);
insert into sources (sou_id, sou_description, sou_idcrm) values (56,'CREDITEA BO',65);
insert into sources (sou_id, sou_description, sou_idcrm) values (57,'SANITAS',66);
insert into sources (sou_id, sou_description, sou_idcrm) values (58,'CREDITEA HM CORTO',67);
insert into sources (sou_id, sou_description, sou_idcrm) values (59,'SANITAS EMISION',68);
insert into sources (sou_id, sou_description, sou_idcrm) values (60,'ABANCA EMISIÓN',69);
insert into sources (sou_id, sou_description, sou_idcrm) values (61,'CREDITEA BBDD ADSALSA',70);
insert into sources (sou_id, sou_description, sou_idcrm) values (62,'CREDITEA MINIPROVI',71);

insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (1,'C2C',2);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (2,'FORM',3);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (3,'ABANDONO',6);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (4,'INACTIVIDAD',7);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (5,'RECICLADO',8);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (6,'PEDIDO NO FINALIZADO LEAD MOVIL',10);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (7,'PEDIDO NO FINALIZADO LEAD SIN CCC',11);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (8,'FDH',12);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (9,'C2C',13);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (10,'PENDIENTE FIRMA',14);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (11,'PEDIDO NO FINALIZADO LEAD COBERTURA',15);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (12,'INICIO PEDIDO',0);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (13,'VENTA',0);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (14,'PENDIENTE FIRMA 2.0',18);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (15,'PENDIENTE EID 2.0',19);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (16,'PENDIENTE CONFIRMA 2.0',20);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (17,'PENDIENTE CAPTACION 2.0',21);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (18,'ENTRANTE',4);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (19,'INCOMPLETOS 2.0',22);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (20,'LL PERDIDA',9);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (21,'DERIVACIÓN YOIGO',23);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (22,'CHAT',24);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (23,'CREDITEA E2E (ENTRANTE-KELISTO)',25);
insert into leadtypes(leatype_id, leatype_description, leatype_idcrm) values (24,'FORM - CONSULTA COBERTURA OK',26);