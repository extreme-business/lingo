interface CookieRepo {
    get(name: string): string | null;
    set(name: string, value: string): void;
    remove(name: string): void;
}

export class CookieRepoImpl implements CookieRepo {
    get(name: string): string | null {
        const cookies = document.cookie.split(';');
        for (const cookie of cookies) {
            const [cookieName, cookieValue] = cookie.split('=');
            if (cookieName.trim() === name) {
                return cookieValue;
            }
        }
        return null;
    }

    set(name: string, value: string): void {
        document.cookie = `${name}=${value}`;
    }

    remove(name: string): void {
        document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC`;
    }
}