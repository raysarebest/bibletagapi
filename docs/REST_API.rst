============
Rest APIs
============

Resource POST ``tag``
---------------------------

REST endpoint for tagging a scripture passage with tag (like a hashtag).  Both book and tag are case insensitive, and book can be a full book name (e.g., "1 Kings" or "Matthew") or a digital bible platform book_id (e.g., "1Kgs" or "Matt").

``http://<host>:<port>/tag``

Request-Response Examples:
~~~~~~~~~~~~~~~~~~~~~~~~~~

Successful case
```````````````

**Request:**

::

    POST /tag

    Content-Type: application/json

    {
    	"tag": "Love",
    	"book": "1Cor",
    	"chapter": 13,
    	"startVerse": 1,
    	"endVerse": 1
    }

**Response:**

::

    HTTP/1.1 200 OK
    Content-Length: 37
    Content-Type: application/json; charset=utf-8
    Accept: application/json

    {
      "code": 200,
      "text": "Tagged Passage"
    }


Error case: could not insert the tag document into the DB
``````````````````````````````````

**Request:**

::

    POST /tag

    Content-Type: application/json

    {
    	"tag": "Love",
    	"book": "1Cor",
    	"chapter": 13,
    	"startVerse": 1,
    	"endVerse": 1
    }

**Response:**

::

    HTTP/1.1 304

Resource GET ``tag/<currenttag>``
---------------------------

REST endpoint for getting a scripture passage tagged with <currenttag>.

``http://<host>:<port>/tag/<currenttag>``

Request-Response Examples:
~~~~~~~~~~~~~~~~~~~~~~~~~~

Successful case
```````````````

**Request:**

::

    GET /tag/love

**Response:**

::

    HTTP/1.1 200 OK
    Content-Length: 285
    Content-Type: application/json; charset=utf-8
    Accept: application/json

    [
      {
        "book_id": "1Cor",
        "book_name": "1 Corinthians",
        "book_order": "61",
        "chapter_id": "13",
        "chapter_title": "Chapter 13",
        "paragraph_number": "140",
        "verse_id": "1",
        "verse_text": "If I speak in the tongues of men and of angels, but have not love, I am a noisy gong or a clanging cymbal. \n\t\t\t"
      }
    ]


Error case: no tagged scripture found
``````````````````````````````````

**Request:**

::

    GET /tag/love

**Response:**

::

    HTTP/1.1 204
