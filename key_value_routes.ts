import {
  Context,
  Router,
  RouterContext,
} from 'https://deno.land/x/oak@v17.1.3/mod.ts';
import {view} from './view_routes.ts';
import {Next} from 'https://deno.land/x/oak@v17.1.3/middleware.ts';
const key_value_router = new Router();

/**A dictionary where the form is
 * {
 *    Example_Key: {
 *            value: 'example_value',
 *            version: 1
 *            },
 *    My_2nd_key: {
 *            value: 'anotherval',
 *            version: 5
 *            }
 * }
 */
export const kv_pairs = {} as {
  [key: string | number]: {version: number; value: any};
};

/**  Lets say metadata has structure
 * {
 *    key: "The Key of the kvp",
 *    version: 1 // The version (number) of the kvp
 *    value: (????) <--- idk if we need this
 * }
 */

async function validate_body(context: any, next: any) {
  const body = await context.request.body.json();
  if (!body) {
    context.response.status = 400;
    console.log('request has no body');
    return;
  }
  if (!Object.hasOwn(body, 'casual-metadata')) {
    context.response.status = 400;
    console.log('request has no casual-metadata attribute');
    return;
  }
  await next();
}

key_value_router.get('/kvs/:key', validate_body, async (context) => {
  // Some code, take note of the bodies casual metadata
  // Check metadata version. If meta-data is GREATER THAN current version, invalid request(?)
  // Get key

  const key = context.params.key;
  const body = await context.request.body.json();

  // Check if key exists
  const value = kv_pairs[key].value;
  if (!value) {
    context.response.status = 404;
    context.response.body = {error: 'Key does not exist'}; //Maybe have to put this inside JSON.stringify()?
    return;
  }

  // Send back found key
  send_res(context, 200, {
    result: 'found',
    value: value,
    'casual-metadata': {key, value, version: 1},
  });
});

key_value_router.put('/kvs/:key', validate_body, async (context) => {
  // Check metadata version, version must be EQUAL or GREATER, if LESS, then reject
  // Some code
  // Update value
  // Propogate updates

  const key = context.params.key;
  const body = await context.request.body.json();

  // Get current version (if any)
  const value = body.value; // Do some error checking on this
  let stored_version = kv_pairs[key]?.version ? kv_pairs[key].version : 0;

  // Check if version from metadata is valid with stored version

  // Store value & increment version

  stored_version++;
  kv_pairs[key] = {value, version: stored_version};

  const casualmetadata = {
    key,
    value,
    version: stored_version,
  };

  let status = 200;
  let result = 'replaced';

  // 1 is initial version, only exists if we just made it
  if (stored_version == 1) {
    status = 201;
    result = 'created';
  }

  context.response.status = status;
  context.response.body = {
    result,
    'casual-metadata': casualmetadata,
  }; //JSON.stringify();
});

key_value_router.delete('/kvs/:key', validate_body, async (context) => {
  // Check metadata version, version must be EQUAL or GREATER, if LESS, then reject
  // Some code
  // Dont delete the entry, simply make the 'value' attribute null/undefined and increment the version.

  const key = context.params.key;
  const body = await context.request.body.json();

  // Get current version (if any)
  let stored_version = kv_pairs[key]?.version ? kv_pairs[key].version : 0;

  // Check if version from metadata is valid with stored version

  // Store value & increment version
  if (stored_version == 0) {
    // return 404;
    return;
  }

  stored_version++;

  // Now get the value
  const stored_value = kv_pairs[key].value;

  // If stored value is undefined, we act like it isnt here
  if (!stored_value) {
    // Return 404
    return;
  }

  // Now we actuall delete it
  kv_pairs[key] = {value: undefined, version: stored_version};

  // return 200 OK
  send_res(context, 200, {
    'casual-metadata': {version: stored_version, key, value: undefined}, // Show metadata that there is currently no value
  });
});

// Replicate values, send the key, value, and version number to all addresses in view
function replicate(key: string, value: any, version: number) {
  const body = {key, value, version};
  view.forEach((v) => {
    // Assuming V is a valid URL
    fetch(`${v}`, {
      method: 'POST',
      body: JSON.stringify(body),
    })
      .then((res) => {
        if (!res.ok) {
          throw res;
        }
      })
      .catch((e) => {
        console.log(`Error ${e}`);
      });
  });
}

function send_res(context: any, status: number, body: any) {
  context.response.status = status;
  context.response.body = body;
}

export default key_value_router;
