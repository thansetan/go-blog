/*
 Navicat Premium Data Transfer

 Source Server         : sanbercode-pgsql
 Source Server Type    : PostgreSQL
 Source Server Version : 150004 (150004)
 Source Host           : localhost:5432
 Source Catalog        : goblog-db
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 150004 (150004)
 File Encoding         : 65001

 Date: 24/08/2023 19:16:59
*/


-- ----------------------------
-- Sequence structure for blogs_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."blogs_id_seq";
CREATE SEQUENCE "public"."blogs_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for comments_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."comments_id_seq";
CREATE SEQUENCE "public"."comments_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for lists_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."lists_id_seq";
CREATE SEQUENCE "public"."lists_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for posts_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."posts_id_seq";
CREATE SEQUENCE "public"."posts_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Sequence structure for users_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."users_id_seq";
CREATE SEQUENCE "public"."users_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 9223372036854775807
START 1
CACHE 1;

-- ----------------------------
-- Table structure for blogs
-- ----------------------------
DROP TABLE IF EXISTS "public"."blogs";
CREATE TABLE "public"."blogs" (
  "id" int8 NOT NULL DEFAULT nextval('blogs_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of blogs
-- ----------------------------
INSERT INTO "public"."blogs" VALUES (1, 'Go-Blog', 'This is the first blog on the platform', 'johndoe', '2023-08-24 11:37:05.865451+00', '2023-08-24 11:40:40.257807+00');
INSERT INTO "public"."blogs" VALUES (2, 'Jimmy Doe''s blog', 'Jimmy Doe''s blog description', 'jimmydoe', '2023-08-24 11:49:44.171977+00', '2023-08-24 11:49:44.171977+00');
INSERT INTO "public"."blogs" VALUES (3, 'Jane Doe''s blog', 'Jane Doe''s blog description', 'janedoe', '2023-08-24 11:54:15.111516+00', '2023-08-24 11:54:15.111516+00');

-- ----------------------------
-- Table structure for comments
-- ----------------------------
DROP TABLE IF EXISTS "public"."comments";
CREATE TABLE "public"."comments" (
  "id" int8 NOT NULL DEFAULT nextval('comments_id_seq'::regclass),
  "commenter" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "post_id" int8 NOT NULL,
  "content" text COLLATE "pg_catalog"."default",
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of comments
-- ----------------------------
INSERT INTO "public"."comments" VALUES (1, 'jimmydoe', 2, 'AYO!? REALLY? I''M SO EXCITED TO HEAR THAT LOL', '2023-08-24 11:53:14.33108+00', '2023-08-24 11:53:14.33108+00');
INSERT INTO "public"."comments" VALUES (2, 'janedoe', 2, 'DISGUSTING, EWWWWWW 🤮🤮🤮🤮🤮', '2023-08-24 11:55:27.420321+00', '2023-08-24 11:55:56.013482+00');
INSERT INTO "public"."comments" VALUES (7, 'janedoe', 3, 'What is this??? 🤪🤪🤪🤪', '2023-08-24 12:15:36.536881+00', '2023-08-24 12:15:36.536881+00');

-- ----------------------------
-- Table structure for list_posts
-- ----------------------------
DROP TABLE IF EXISTS "public"."list_posts";
CREATE TABLE "public"."list_posts" (
  "list_id" int8,
  "post_id" int8
)
;

-- ----------------------------
-- Records of list_posts
-- ----------------------------
INSERT INTO "public"."list_posts" VALUES (1, 2);
INSERT INTO "public"."list_posts" VALUES (3, 2);
INSERT INTO "public"."list_posts" VALUES (3, 1);
INSERT INTO "public"."list_posts" VALUES (4, 1);

-- ----------------------------
-- Table structure for lists
-- ----------------------------
DROP TABLE IF EXISTS "public"."lists";
CREATE TABLE "public"."lists" (
  "id" int8 NOT NULL DEFAULT nextval('lists_id_seq'::regclass),
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "slug" varchar(510) COLLATE "pg_catalog"."default" NOT NULL,
  "description" text COLLATE "pg_catalog"."default",
  "owner" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of lists
-- ----------------------------
INSERT INTO "public"."lists" VALUES (1, 'My Favourite 🤩🤩🤩🤩', 'my-favourite-1692878357', 'I love these things', 'jimmydoe', '2023-08-24 11:59:17.2951+00', '2023-08-24 11:59:17.2951+00');
INSERT INTO "public"."lists" VALUES (3, 'My Favourite 🤩🤩🤩🤩', 'my-favourite-1692878859', 'I love these things', 'jimmydoe', '2023-08-24 12:07:40.077819+00', '2023-08-24 12:07:40.077819+00');
INSERT INTO "public"."lists" VALUES (4, 'interesting things', 'interesting-things-1692879040', 'posts I found interesting', 'janedoe', '2023-08-24 12:10:40.245575+00', '2023-08-24 12:10:40.245575+00');

-- ----------------------------
-- Table structure for posts
-- ----------------------------
DROP TABLE IF EXISTS "public"."posts";
CREATE TABLE "public"."posts" (
  "id" int8 NOT NULL DEFAULT nextval('posts_id_seq'::regclass),
  "title" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "slug" varchar(510) COLLATE "pg_catalog"."default" NOT NULL,
  "content" text COLLATE "pg_catalog"."default" NOT NULL,
  "blog_id" int8 NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of posts
-- ----------------------------
INSERT INTO "public"."posts" VALUES (1, 'First Post Let''s GOOOOOO', 'first-post-let-s-goooooo-1692877400', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed dapibus ligula et neque gravida faucibus. Sed elit neque, tincidunt eget nibh at, mollis scelerisque mi. Sed non libero sit amet nisi porttitor egestas vel eget magna. Praesent in eros elementum elit mattis sollicitudin. Mauris et tristique est. Mauris ornare tempor eros. Duis semper tincidunt venenatis. Pellentesque eleifend id tellus eu sollicitudin. Donec fringilla nunc at malesuada dignissim. Integer a dictum ante. Proin eu vehicula felis. Quisque consequat felis dolor. Nulla ut eleifend eros, at dictum ante. Sed condimentum lectus vitae lectus dignissim mollis. Integer congue neque dolor, ac euismod mauris eleifend ac. Donec condimentum consequat lacus non rhoncus. Phasellus posuere neque lectus. Proin sollicitudin, neque sed vulputate vestibulum, neque ipsum congue neque, vitae mattis velit eros eu enim. Donec lectus ante, tristique at molestie eu, placerat sit amet ex. Duis at placerat eros, vitae commodo diam. Suspendisse ut sollicitudin velit. Vestibulum et iaculis arcu. Aenean sit amet elit quis lectus porttitor ornare sed quis ipsum. Donec ac congue nibh. Cras vulputate lectus sit amet viverra dapibus.', 1, '2023-08-24 11:43:20.775158+00', '2023-08-24 11:43:20.775158+00');
INSERT INTO "public"."posts" VALUES (2, 'DID YOU KNOW?', 'did-you-know-1692877455', 'Hey guys, did you know that in terms of male human and female Pokémon breeding, Vaporeon is the most compatible Pokémon for humans? Not only are they in the field egg group, which is mostly comprised of mammals, Vaporeon are an average of 3”03’ tall and 63.9 pounds, this means they’re large enough to be able handle human dicks, and with their impressive Base Stats for HP and access to Acid Armor, you can be rough with one. Due to their mostly water based biology, there’s no doubt in my mind that an aroused Vaporeon would be incredibly wet, so wet that you could easily have sex with one for hours without getting sore. They can also learn the moves Attract, Baby-Doll Eyes, Captivate, Charm, and Tail Whip, along with not having fur to hide nipples, so it’d be incredibly easy for one to get you in the mood. With their abilities Water Absorb and Hydration, they can easily recover from fatigue with enough water. No other Pokémon comes close to this level of compatibility. Also, fun fact, if you pull out enough, you can make your Vaporeon turn white. Vaporeon is literally built for human dick. Ungodly defense stat+high HP pool+Acid Armor means it can take cock all day, all shapes and sizes and still come for more', 1, '2023-08-24 11:44:15.196721+00', '2023-08-24 11:44:15.196721+00');
INSERT INTO "public"."posts" VALUES (3, 'this is good', 'this-is-good-1692877563', '👌👀👌👀👌👀👌👀👌👀 good shit go౦ԁ sHit👌 thats ✔ some good👌👌shit right👌👌th 👌 ere👌👌👌 right✔there ✔✔if i do ƽaү so my selｆ 💯 i say so 💯 thats what im talking about right there right there (chorus: ʳᶦᵍʰᵗ ᵗʰᵉʳᵉ) mMMMMᎷМ💯 👌👌 👌НO0ОଠＯOOＯOОଠଠOoooᵒᵒᵒᵒᵒᵒᵒᵒᵒ👌 👌👌 👌 💯 👌 👀 👀 👀 👌👌Good shit', 1, '2023-08-24 11:46:03.99193+00', '2023-08-24 11:46:55.963668+00');

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS "public"."users";
CREATE TABLE "public"."users" (
  "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
  "email" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT NULL::character varying,
  "username" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL DEFAULT NULL::character varying,
  "password" bytea NOT NULL,
  "created_at" timestamptz(6),
  "updated_at" timestamptz(6)
)
;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO "public"."users" VALUES (1, 'johndoe@gmail.com', 'johndoe', 'John Doe', E'$2a$05$BvYYWo9qMCrjU5VyG6RDaeB/D3KHd/qkEoAFJRnpbxOtFy9QtsOUS', '2023-08-24 11:37:05.861415+00', '2023-08-24 11:37:05.861415+00');
INSERT INTO "public"."users" VALUES (4, 'janedoe@gmail.com', 'janedoe', 'Jane Doe', E'$2a$05$6VPqxju7/9t9TzS1RmaLfOccrZ/hgcvB5sSNKEQXyL0tmGQyd98Zu', '2023-08-24 11:54:15.109907+00', '2023-08-24 11:54:15.109907+00');
INSERT INTO "public"."users" VALUES (3, 'jimmydoe.2@gmail.com', 'jimmydoe', 'Jimmy Doe The II', E'$2a$05$/7xXQXV5wcbE7krDNYO35e3S0SyQlJ.7fLKn3DgxzP11k4Iha8uLq', '2023-08-24 11:49:44.16867+00', '2023-08-24 12:12:30.660773+00');

-- ----------------------------
-- Function structure for uuid_generate_v1
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v1mc
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v1mc"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v1mc"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v1mc'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v3
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v3"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v3"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v3'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v4
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v4"();
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v4"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v4'
  LANGUAGE c VOLATILE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_generate_v5
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_generate_v5"("namespace" uuid, "name" text);
CREATE OR REPLACE FUNCTION "public"."uuid_generate_v5"("namespace" uuid, "name" text)
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_generate_v5'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_nil
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_nil"();
CREATE OR REPLACE FUNCTION "public"."uuid_nil"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_nil'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_dns
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_dns"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_dns"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_dns'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_oid
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_oid"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_oid"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_oid'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_url
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_url"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_url"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_url'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Function structure for uuid_ns_x500
-- ----------------------------
DROP FUNCTION IF EXISTS "public"."uuid_ns_x500"();
CREATE OR REPLACE FUNCTION "public"."uuid_ns_x500"()
  RETURNS "pg_catalog"."uuid" AS '$libdir/uuid-ossp', 'uuid_ns_x500'
  LANGUAGE c IMMUTABLE STRICT
  COST 1;

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."blogs_id_seq"
OWNED BY "public"."blogs"."id";
SELECT setval('"public"."blogs_id_seq"', 3, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."comments_id_seq"
OWNED BY "public"."comments"."id";
SELECT setval('"public"."comments_id_seq"', 7, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."lists_id_seq"
OWNED BY "public"."lists"."id";
SELECT setval('"public"."lists_id_seq"', 4, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."posts_id_seq"
OWNED BY "public"."posts"."id";
SELECT setval('"public"."posts_id_seq"', 5, true);

-- ----------------------------
-- Alter sequences owned by
-- ----------------------------
ALTER SEQUENCE "public"."users_id_seq"
OWNED BY "public"."users"."id";
SELECT setval('"public"."users_id_seq"', 4, true);

-- ----------------------------
-- Indexes structure for table blogs
-- ----------------------------
CREATE UNIQUE INDEX "idx_blogs_owner" ON "public"."blogs" USING btree (
  "owner" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table blogs
-- ----------------------------
ALTER TABLE "public"."blogs" ADD CONSTRAINT "blogs_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Primary Key structure for table comments
-- ----------------------------
ALTER TABLE "public"."comments" ADD CONSTRAINT "comments_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table list_posts
-- ----------------------------
CREATE UNIQUE INDEX "idx_list_id_post_id" ON "public"."list_posts" USING btree (
  "list_id" "pg_catalog"."int8_ops" ASC NULLS LAST,
  "post_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Indexes structure for table lists
-- ----------------------------
CREATE INDEX "idx_lists_slug" ON "public"."lists" USING btree (
  "slug" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table lists
-- ----------------------------
ALTER TABLE "public"."lists" ADD CONSTRAINT "lists_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table posts
-- ----------------------------
CREATE INDEX "idx_posts_blog_id" ON "public"."posts" USING btree (
  "blog_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);
CREATE INDEX "idx_posts_slug" ON "public"."posts" USING btree (
  "slug" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table posts
-- ----------------------------
ALTER TABLE "public"."posts" ADD CONSTRAINT "posts_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table users
-- ----------------------------
CREATE UNIQUE INDEX "idx_users_username" ON "public"."users" USING btree (
  "username" COLLATE "pg_catalog"."default" "pg_catalog"."text_ops" ASC NULLS LAST
);

-- ----------------------------
-- Uniques structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_email_key" UNIQUE ("email");

-- ----------------------------
-- Primary Key structure for table users
-- ----------------------------
ALTER TABLE "public"."users" ADD CONSTRAINT "users_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Foreign Keys structure for table blogs
-- ----------------------------
ALTER TABLE "public"."blogs" ADD CONSTRAINT "fk_blogs_user" FOREIGN KEY ("owner") REFERENCES "public"."users" ("username") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table comments
-- ----------------------------
ALTER TABLE "public"."comments" ADD CONSTRAINT "fk_comments_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON DELETE CASCADE ON UPDATE NO ACTION;
ALTER TABLE "public"."comments" ADD CONSTRAINT "fk_comments_user" FOREIGN KEY ("commenter") REFERENCES "public"."users" ("username") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table list_posts
-- ----------------------------
ALTER TABLE "public"."list_posts" ADD CONSTRAINT "fk_list_posts_post" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON DELETE CASCADE ON UPDATE CASCADE;
ALTER TABLE "public"."list_posts" ADD CONSTRAINT "fk_lists_list_posts" FOREIGN KEY ("list_id") REFERENCES "public"."lists" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table lists
-- ----------------------------
ALTER TABLE "public"."lists" ADD CONSTRAINT "fk_lists_user" FOREIGN KEY ("owner") REFERENCES "public"."users" ("username") ON DELETE NO ACTION ON UPDATE NO ACTION;

-- ----------------------------
-- Foreign Keys structure for table posts
-- ----------------------------
ALTER TABLE "public"."posts" ADD CONSTRAINT "fk_blogs_posts" FOREIGN KEY ("blog_id") REFERENCES "public"."blogs" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION;
