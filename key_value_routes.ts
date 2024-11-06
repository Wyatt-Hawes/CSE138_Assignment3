import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';

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
 * }
 */
key_value_router.get('/kvs/:value', (context) => {
  // Some code, take note of the bodies casual metadata

  // Check metadata version. If meta-data is GREATER THAN current version, invalid request(?)

  const value = context.params.value;
  context.response.body = `Hello World! You entered: ${value}`;
  context.response.status = 201;
});

key_value_router.put('/kvs/:value', (context) => {
  // Check metadata version, version must be EQUAL or GREATER, if LESS, then reject
  // Some code
  // Update value
  // Propogate updates
});

key_value_router.delete('/kvs/:values', (context) => {
  // Check metadata version, version must be EQUAL or GREATER, if LESS, then reject
  // Some code
  // Dont delete the entry, simply make the 'value' attribute null/undefined and increment the version.
});

export default key_value_router;
