FROM scratch
ADD files /files
ADD bibletagapi /bibletagapi
EXPOSE 8080
CMD ["/bibletagapi"]