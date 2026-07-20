# Question Module API

Source: [Feishu - Interface Document (2): Question Module](https://acn2lw4rwc26.feishu.cn/wiki/KXGdwKDKGiCfYYkgI6OcU3Vfndd)

Snapshot date: 2026-07-21. This snapshot records the sections verified for the
current implementation; consult Feishu before changing an unrecorded endpoint.

## Common Response

```json
{
  "code": 0,
  "msg": "success",
  "data": {}
}
```

Question query endpoints use soft authentication. The current backend has no
request identity middleware, so list responses follow the visitor behavior.

## GET /api/v1/question/all/list

Optional query: `job_name`.

The response data is grouped by `job_name`, then contains `series_list` and each
series contains `bank_list`. Question counts include only normal questions
(`question.status = 1`). Series and banks are ordered by `order` ascending.

## GET /api/v1/question/list

Optional query parameters:

| Name | Type | Meaning |
| --- | --- | --- |
| `bank_id` | string | Question bank membership |
| `keyword` | string | Case-insensitive match against `title` or `content` |
| `difficulty` | int | Difficulty filter |
| `tags` | string[] | All supplied tags must match (`$all`) |
| `job_name` | string | Job direction filter |
| `user_tag` | int | User-specific mastery filter; requires authenticated user context |
| `sort_by` | string | `create_time`, `view_count`, `thumbs_up_count`, or `dislike_count` |
| `sort_order` | string | `asc` or `desc` |
| `page` | int | Default `1` |
| `page_size` | int | Default `20` |

Current visitor defaults:

- only normal questions (`status = 1`);
- `sort_by = order`, `sort_order = asc`;
- `user_tag = 0`, `user_liked = false` in every response item;
- `content` is truncated to the first 150 Unicode characters.

Response data:

```json
{
  "total": 1,
  "list": [
    {
      "question_id": "q001",
      "bank_list": ["b001"],
      "title": "Question title",
      "content": "Question content",
      "difficulty": 1,
      "tags": ["Go"],
      "status": 1,
      "vip": false,
      "hot_degree": 0,
      "view_count": 0,
      "thumbs_up_count": 0,
      "dislike_count": 0,
      "order": 1,
      "user_tag": 0,
      "user_liked": false,
      "create_time": "2026-07-21 00:00:00",
      "update_time": "2026-07-21 00:00:00"
    }
  ]
}
```

## GET /api/v1/question/meta/list

Uses the same query parameters, filtering, sorting, and pagination behavior as
`GET /api/v1/question/list`.

Response data contains `total` and `list`. Each list item only contains
`question_id` and `title`; it is intended for the question navigation sidebar.

```json
{
  "total": 100,
  "list": [
    {
      "question_id": "q001",
      "title": "Question title"
    }
  ]
}
```

## GET /api/v1/question/detail

Required query parameter: `question_id`.

The response `data` is one complete question using the same field contract as a
`GET /api/v1/question/list` item. At the current implementation stage, content
is returned in full. `user_tag` and `user_liked` use the same visitor rules as
the list endpoint until authentication and user interaction lookup are added.

Only normal questions (`status = 1`) are visible. A missing question returns the
standard response envelope with business code `404` and `data: null`.

## GET /api/v1/question/hot/list

Optional query parameters:

- `limit`: maximum number of questions, default `10`.
- `job_name`: job direction filter.

Only normal questions (`status = 1`) are returned. Results are always sorted by
`hot_degree` descending; `view_count` is a response field, not the ranking key.
Each item contains only `question_id`, `bank_list`, `title`, and `view_count`.

```json
{
  "list": [
    {
      "question_id": "q001",
      "bank_list": ["b001"],
      "title": "Question title",
      "view_count": 1024
    }
  ]
}
```

## Other Confirmed Endpoints

- `POST /api/v1/safe/tag_question`
