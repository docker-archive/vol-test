# volcheck
#
#
# This Dockerfile expects a volume mounted on /data
# It will write a couple of files at startup in that volume
# and make the status of those files available as http calls
# for verification of volume mounts.
FROM alpine:3.8

LABEL maintainer="keith@docker.com"

COPY textcheck.shtml /var/www/localhost/htdocs/textcheck.shtml
COPY bincheck.shtml /var/www/localhost/htdocs/bincheck.shtml
COPY runfilecheck.shtml /var/www/localhost/htdocs/runfilecheck.shtml
COPY idempotent_filecheck.sh /usr/bin/idempotent_filecheck.sh
RUN apk add --no-cache apache2 apache2-utils && \
    mkdir -p /run/apache2/ && \
    cp /etc/apache2/httpd.conf /etc/apache2/httpd.orig && \
    echo "Options +Includes" > /var/www/.htaccess # && \
    chmod +x /var/www/localhost/htdocs/textcheck.shtml
ADD httpd.conf /etc/apache2/httpd.conf
RUN chmod +x /usr/bin/idempotent_filecheck.sh

EXPOSE 80

CMD /usr/sbin/httpd -f /etc/apache2/httpd.conf -DFOREGROUND
