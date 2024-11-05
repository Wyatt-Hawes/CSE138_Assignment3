import {Router} from 'https://deno.land/x/oak@v17.1.3/mod.ts';

// All these calls will directly update values but not replicate any changes

const internal_router = new Router();

// All routes call same functions but without causing a replication
internal_router.get('/kvs/:value', (context) => {
  // Some code
});

internal_router.put('/kvs/:value', (context) => {
  // Some code
});

internal_router.delete('/kvs/:values', (context) => {
  // Some code
});

internal_router.get('/view', (context) => {
  // Some code
});

internal_router.put('/view', (context) => {
  // Some code
});

internal_router.delete('/view', (context) => {
  // Some code
});

export default internal_router;
