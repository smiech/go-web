import { Component, OnInit, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { AuthService } from '../../services/auth.service';
import { PersistenceService } from '../../services/persistence.service';
import { AuthRequestResult } from '../../models/request-result';
import { Router } from '@angular/router';
import { ActivatedRoute } from '@angular/router';

@Component({
    selector: 'login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit, AfterViewInit {
    private userName: string;
    private password: string;
    private showLoginError: boolean = false;
    public showPassword: boolean = false;
    @ViewChild('loginInput') loginInput: ElementRef;

    constructor(private authService: AuthService, private router: Router, private route: ActivatedRoute,
        private persistenceService: PersistenceService) {
    }

    public ngOnInit(): void {
        if (this.authService.checkLogin()) {
            this.router.navigate(['/userspace']);
        }

        this.route.params.subscribe(params => {
            this.userName = params['email'];
        });

    }

    public ngAfterViewInit(): void {
        if (this.loginInput.nativeElement && this.loginInput.nativeElement.focus) {
            this.loginInput.nativeElement.focus();
        }
    }

    public login(): void {
        this.authService.login(this.userName, this.password)
            .then(x => this.handleAuth(x));
    }

    public checklogin(): void {
        this.authService.getUserInfo()
            .then(x => console.log(x));
    }

    private handleAuth(user: AuthRequestResult): void {
        if (user.State < 1) {
            this.showLoginError = true;
        } else {
            this.showLoginError = false;
            if (this.router) {

                this.router.navigate(['/userspace']);
            };
        }
    }
}
