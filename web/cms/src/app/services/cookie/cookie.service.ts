import { Injectable } from '@angular/core';
import { CookieRepoImpl } from './repo';

@Injectable({
  providedIn: 'root'
})
export class CookieService {
  private cookieRepo: CookieRepoImpl;

  constructor() { 
    this.cookieRepo = new CookieRepoImpl();
  }

  get(name: string): string | null {
    return this.cookieRepo.get(name);
  }

  set(name: string, value: string): void {
    this.cookieRepo.set(name, value);
  }

  remove(name: string): void {
    this.cookieRepo.remove(name);
  }
}
