# Comment and Interaction API

Source: [Feishu - Comment and Interaction API](https://acn2lw4rwc26.feishu.cn/wiki/UYN6w62h9ixNlskeeHUccrmFnDh)

Snapshot date: 2026-07-22. Field names below are API contracts.

## GET /api/v1/open/list_comments

Authentication is soft. Missing or invalid credentials are treated as a
visitor; a valid JWT supplies the `user_id` used to derive `user_liked`.

Query parameters:

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `target_type` | int | yes | `1` question / `2` interview experience |
| `target_id` | string | yes | Target ID |
| `parent_id` | string | no | Empty lists top-level comments; otherwise lists that comment's children |
| `sort_by` | string | no | `create_time` or `thumbs_up`; defaults to `create_time` |
| `sort_order` | string | no | `asc` or `desc`; defaults to `desc` |
| `page` | int | no | Defaults to `1` |
| `page_size` | int | no | Defaults to `20` |
| `sub_comment_page` | int | no | Per-parent child page; defaults to `1` |
| `sub_comment_size` | int | no | Per-parent child page size; defaults to `5` |

Successful response data contains `total` and `list`. Each comment contains:

- `comment_id`, `user_id`, `user_name`, `user_avatar`, `content`
- `parent_id`, `reply_to`, `reply_to_name`, `status`, `thumbs_up`
- `sub_comment_total`, `user_liked`, `sub_comments`
- `create_time`, `update_time`

Top-level comments include the requested page of second-level comments in
`sub_comments`. When `parent_id` is supplied, the response is a flat child
comment list and does not nest further `sub_comments`.

## POST /api/v1/comment/add

Authentication is required. The JWT middleware supplies `user_id`.

Request body:

| Field | Type | Required | Notes |
| --- | --- | --- | --- |
| `target_type` | int | yes | `1` question / `2` interview experience |
| `target_id` | string | yes | Target ID |
| `parent_id` | string | no | Empty for a top-level comment |
| `reply_to` | string | no | User ID being replied to |
| `content` | string | yes | Stored after sensitive-word filtering |

The service generates a UUID `comment_id`, stores the comment with `status = 2`,
and increments the parent comment's `child_count` for a reply. Successful
response data contains `comment_id` and the complete `comment` object with the
author name and avatar.
