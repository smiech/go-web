import { User } from '../models/auth-models';
import { Injectable } from '@angular/core';


@Injectable()
export class PersistenceService {
    get accessToken(): string {
        return this.getItem('token');
    }

    set accessToken(token: string) {
        this.setItem('token', token);
    }

    get currentUser(): User {
        return this.getItem('currentUser');
    }

    set currentUser(user: User) {
        this.setItem('currentUser', user);
    }

    get planId(): number {
        return this.getItem('planId');
    }

    set planId(value: number) {
        this.setItem('planId', value);
    }

    get isFirstLogin(): boolean {
        let value = this.getItem('isFirstLogin');

        return value === undefined || value === null ? false : value;
    }

    set isFirstLogin(value: boolean) {
        this.setItem('isFirstLogin', value);
    }

    get showCompltetedMeasures(): boolean {
        let value = this.getItem('showCompltetedMeasures');
        return value === undefined || value === null ? false : value;
    }

    set showCompltetedMeasures(value: boolean) {
        this.setItem('showCompltetedMeasures', value);
    }

    get showAttachedFiles(): boolean {
        let value = this.getItem('showAttachedFiles');
        return value === undefined || value === null ? false : value;
    }

    set showAttachedFiles(value: boolean) {
        this.setItem('showAttachedFiles', value);
    }

    /* --- private area --- */
    private getItem(key: string): any {
        if (typeof window !== 'undefined') {
            let rawValue = localStorage.getItem(key);
            return JSON.parse(rawValue);
        }

        return null;
    }

    private setItem(key: string, value: any): any {
        if (typeof window !== 'undefined') {
            if (value == null)
            {
                localStorage.removeItem(key);
            } else {
                let rawValue = JSON.stringify(value);
                localStorage.setItem(key, rawValue);
            }            
        }
        return null;
    }
}