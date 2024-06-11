import { CanActivateFn } from '@angular/router';

export const accountGuard: CanActivateFn = (route, state) => {
  // check for cookies
  

  return true;
};
