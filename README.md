# qgen

Генератор отчетов СОРМ для биллинга abills

### usage:

```
QGEN_MYSQL="abills:superpassword@/abills?parseTime=true" \
QGEN_PATH=tmp \
QGEN_REGION_ID="42" \
QGEN_REGION_NAME="ВоблаТелком" \
QGEN_INIT_DATE="2012-02-05" \
QGEN_ONLY_ONE_DAY=false \
./qgen
```

SQL create log table

```sql
CREATE TABLE qgenlog (
  id int(11) unsigned NOT NULL AUTO_INCREMENT,
  date datetime DEFAULT NULL,
  comment varchar(512) DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```

for old freebsd system compile gover 1.16

```
GOOS=freebsd GOARCH="amd64"  /opt/homebrew/Cellar/go@1.16/1.16.15/bin/go build
```