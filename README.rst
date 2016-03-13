====================
BibleTag-API
====================

| This service supports a REST API for collaboratively tagging the Bible.

`For more information on the endpoints, visit the API docs. <docs>`_

----
Data
----

This service maintains an authority on arbitrary, crowdsourced tags (e.g., "love", "encouragement", or "frustration") and how those tags are connected to scripture passages (e.g., 1 Cor. 13:1 or John 3:16).

------------
Technologies
------------

* `Go <https://golang.org/>`_
* `RethinkDB <https://www.rethinkdb.com/>`_
* `Docker <https://www.docker.com/>`_ - `DockerFile <DockerFile>`_

-----
Links
-----

* `Unit Tests <goserver_test.go>`_
* `Example Environmental Vars <files/example.env>`_

--------
License
--------

BibleTag-API was created by Daniel Whitenack, Michael Hulet, Josh Oppenheim, Benjamin Bledsoe, and anthonyt2345 and is licensed under an `MIT-style License <License.md>`_.
