import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'lg-button',
  templateUrl: './lg-button.component.html',
  styleUrls: ['./lg-button.component.scss']
})
export class LgButtonComponent implements OnInit {
  @Input() LgButtonClass:string;
  constructor() { }

  ngOnInit() {
  }

}
