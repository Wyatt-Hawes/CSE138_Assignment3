import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';

// All these calls will directly update values but not replicate any changes

const internal_router = new Router();

// Lets say this is to be called by other systems to notify us of any updates regarding our data
internal_router.get('/update', (context) => {
  // Some code
});

export default internal_router;
