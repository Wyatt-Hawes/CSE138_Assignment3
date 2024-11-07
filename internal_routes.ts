import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';
import {kv_pairs} from './key_value_routes.ts';
// All these calls will directly update values but not replicate any changes

const internal_router = new Router();

// Lets say this is to be called by other systems to notify us of any updates regarding our data
internal_router.post('/update', (context) => {
  // Some code
  // Reflect change in kv_pairs
});

export default internal_router;
