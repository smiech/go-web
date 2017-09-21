import { Injectable } from '@angular/core';

@Injectable()
export class ConfigService {
    public serviceUrl: string = 'localhost:5000';
    public getAbsuluteUrl(relativeUrl: string): string {
        return 'http://localhost:50000' + relativeUrl;
    }

    public getAbsoluteSecureUrl(relativeUrl: string) {
        return 'https://localhost:44381' + relativeUrl;
    }
}