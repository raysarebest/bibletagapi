FROM scratch
ADD bibletagapi /bibletagapi
EXPOSE 8080
CMD ["/bibletagapi"]