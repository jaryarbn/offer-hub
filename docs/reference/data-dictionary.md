# Offer Hub Data Dictionary (Verified Sections)

Source: [Feishu - Data Dictionary](https://acn2lw4rwc26.feishu.cn/wiki/VEfQwFVh7iaLAZk1iLGcNsMJnUZ)

`analysis_content` is added by the verified question-detail authentication
lesson contract. The remaining schemas below come from the Feishu data
dictionary.

Snapshot date: 2026-07-22. MongoDB field names below are persistence contracts.

## question

| Field | Type | Notes |
| --- | --- | --- |
| `_id` | ObjectID | Mongo primary key |
| `question_id` | string | Business ID |
| `bank_list` | string[] | Bank IDs |
| `job_name` | string | Job direction |
| `title` | string | Title |
| `content` | string | Content |
| `analysis_content` | string | Answer analysis; complete value is exposed only to authenticated detail requests |
| `difficulty` | int | Difficulty |
| `tags` | string[] | Tags |
| `status` | int | `1` means normal |
| `vip` | bool | VIP-only flag |
| `hot_degree` | int | Hot score |
| `view_count` | int | View count |
| `thumbs_up_count` | int | Like count |
| `dislike_count` | int | Dislike count |
| `order` | int64 | Manual order |
| `create_time` | datetime | Creation time |
| `update_time` | datetime | Update time |

`user_tag` and `user_liked` are response fields derived from user interaction
data. They are not fields of a question document.

## comments

| Field | Type | Notes |
| --- | --- | --- |
| `_id` | string | Mongo primary key |
| `comment_id` | string | Business comment ID |
| `target_type` | int | `1` question / `2` interview experience |
| `target_id` | string | Comment target ID |
| `question_id` | string | Associated question ID |
| `user_id` | string | Comment author ID |
| `content` | string | Comment content |
| `parent_id` | string | Empty for a top-level comment; otherwise the parent comment ID |
| `reply_to` | string | User ID being replied to; empty for a top-level comment |
| `child_count` | int | Stored child count; second-level comments use `0` |
| `thumbs_up` | int | Like count |
| `view_count` | int | View count |
| `status` | int | `1` reviewing / `2` normal / `3` rejected / `4` hidden / `5` deleted |
| `create_time` | datetime | Creation time |
| `update_time` | datetime | Update time |

The list API derives `user_name`, `user_avatar`, `reply_to_name`,
`sub_comment_total`, `sub_comments`, and `user_liked`; these are not persisted
fields in `comments`.

## user_interactions

| Field | Type | Notes |
| --- | --- | --- |
| `_id` | ObjectID | Mongo primary key |
| `user_id` | string | User ID |
| `target_type` | int | `1` question / `2` interview experience / `3` comment |
| `target_id` | string | Target business ID |
| `interaction_type` | int | `1` like / `2` dislike |
| `status` | int | `1` active / `0` cancelled |
| `create_time` | datetime | Creation time |
| `update_time` | datetime | Update time |

For question detail, `user_liked` is true only when `user_id`,
`target_type = 1`, `target_id = question_id`, `interaction_type = 1`, and
`status = 1` all match.

## user_question_tag

| Field | Type | Notes |
| --- | --- | --- |
| `_id` | ObjectID | Mongo primary key |
| `user_id` | string | User ID |
| `question_id` | string | Question ID |
| `tag` | int | `0` unmarked / `1` mastered / `2` review later / `3` not mastered |
| `create_time` | datetime | Creation time |
| `update_time` | datetime | Last update time |

## question_bank_series

| Field | Type |
| --- | --- |
| `_id` | ObjectID |
| `series_id` | string |
| `series_name` | string |
| `job_name` | string |
| `order` | int64 |
| `create_time` | datetime |
| `update_time` | datetime |

## question_bank

| Field | Type |
| --- | --- |
| `_id` | ObjectID |
| `bank_id` | string |
| `series_id` | string |
| `bank_name` | string |
| `bank_logo` | string |
| `desc` | string |
| `job_name` | string |
| `order` | int64 |
| `create_time` | datetime |
| `update_time` | datetime |

The API `count` field is aggregated from `question`; it is not stored in
`question_bank`.
