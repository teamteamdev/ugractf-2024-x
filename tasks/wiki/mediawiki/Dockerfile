FROM mediawiki:1.33
RUN apt-get update && apt-get -y install sqlite3

COPY logo.svg /var/www/html/
COPY LocalSettings.php /var/www/html/
COPY my_wiki.sqlite.sql /var/www/data/
RUN sqlite3 /var/www/data/my_wiki.sqlite </var/www/data/my_wiki.sqlite.sql \
    && chown www-data:www-data /var/www/data/my_wiki.sqlite \
    && rm /var/www/data/my_wiki.sqlite.sql
