import { WiContrib, WiContributionHandler } from
  "wi-studio/app/contrib/wi-contrib";
import { IActivityContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult } from
  "wi-studio/common/models/validation";
import { Injectable, Inject, Injector } from "@angular/core";
import { Observable } from "rxjs/Observable";
@WiContrib({
  handler: ConcatActivityUIContributionHandler
})
@Injectable()
export class ConcatActivityUIContributionHandler extends
  WiContributionHandler {
  constructor( @Inject(Injector) injector) {
    super(injector);
  }
  getValue(fieldName: string, context: IActivityContribution):
    Observable<any> | any {
    if (fieldName === "separator") {
      let list: Array<string> = ["-", "$", "#"];
      return list;
    }
    return null;
  }
  validate(fieldName: string, context: IActivityContribution):
    Observable<IValidationResult> | IValidationResult {
    if (fieldName === "separator") {
      let vresult: IValidationResult =
        ValidationResult.newValidationResult();
      if (context.getField("useSeparator").value &&
        context.getField("useSeparator").value === true) {
        if (context.getField("separator").display &&
          context.getField("separator").display.visible == false {
          vresult.setVisible(true);
        }
        if (context.getField("separator").value === null ||
          context.getField("separator").value === "") {
          vresult.setError("Separator must be configured");
        }
      } else {
        vresult.setVisible(false);
      }
      return vresult;
    }
    return null;
  }
}
