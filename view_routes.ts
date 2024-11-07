import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';
const view_router = new Router();
export const view = [];
// Get VIEW from environmental variable

view_router.get('/', (context) => {
  context.response.body = 'Hello World!';
  context.response.status = 201;
});

view_router.put('/view', (context) => {
  const body = context.request.body;
  // Add to view
});

view_router.get('/view', (context) => {
  // Return view
});

view_router.delete('/view', (context) => {
  // Delete from view
});

export default view_router;
