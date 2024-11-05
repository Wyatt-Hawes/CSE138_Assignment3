import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';

const key_value_router = new Router();

key_value_router.get('/kvs/:value', (context) => {
  // Some code
  const value = context.params.value;
  context.response.body = `Hello World! You entered: ${value}`;
  context.response.status = 201;
});

key_value_router.put('/kvs/:value', (context) => {
  // Some code
});

key_value_router.delete('/kvs/:values', (context) => {
  // Some code
});

export default key_value_router;
