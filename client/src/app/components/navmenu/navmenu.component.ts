import { Component } from '@angular/core';
import { Menuitem } from '../../models/menuitem-models';
@Component({
    selector: 'nav-menu',
    templateUrl: './navmenu.component.html',
    styleUrls: ['./navmenu.component.css']
})
export class NavMenuComponent {
    public extraMenuItems: Array<Menuitem> = [{
        name:'extra',
        link:'extra'
    }];
    updateMenuItems(value) {
        this.extraMenuItems = value;
      }
}
