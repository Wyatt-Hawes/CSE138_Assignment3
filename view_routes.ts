import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';
import {my_obj} from './key_value_routes.ts';
const view_router = new Router();

view_router.get('/', (context) => {
  context.response.body = 'Hello World!';
  context.response.status = 201;
});

view_router.put('/view', (context) => {
  const body = context.request.body;
  // Some code
});

view_router.get('/view', (context) => {
  // Some code
});

view_router.delete('/view', (context) => {
  // Some code
});

export default view_router;
