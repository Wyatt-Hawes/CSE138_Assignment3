import {Application} from 'https://deno.land/x/oak@v17.1.3/mod.ts';
import view_router from './view_routes.ts';
import key_value_router from './key_value_routes.ts';

const app = new Application();
const port = 8090;

app.use(view_router.routes());
app.use(key_value_router.routes());

app.use(view_router.allowedMethods());
app.use(key_value_router.allowedMethods());

console.log('🦕 Server is now Listening 🦕');
await app.listen({port});
