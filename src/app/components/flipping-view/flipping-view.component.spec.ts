import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FlippingViewComponent } from './flipping-view.component';

describe('FlippingViewComponent', () => {
  let component: FlippingViewComponent;
  let fixture: ComponentFixture<FlippingViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FlippingViewComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(FlippingViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
