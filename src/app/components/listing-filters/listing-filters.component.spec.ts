import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListingFiltersComponent } from './listing-filters.component';

describe('ListingFiltersComponent', () => {
  let component: ListingFiltersComponent;
  let fixture: ComponentFixture<ListingFiltersComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListingFiltersComponent]
    })
    .compileComponents();
    
    fixture = TestBed.createComponent(ListingFiltersComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
