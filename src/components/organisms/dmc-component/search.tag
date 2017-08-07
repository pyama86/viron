dmc-component-search.ComponentSearch
  .ComponentSearch__head
    .ComponentSearch__name { opts.name }
    .ComponentSearch__description { opts.description }
  .ComponentSearch__body
    dmc-parameters(parameterObjects="{ opts.parameterObjects }" initialParameters="{ opts.initialparameters }" onChange="{ handleParametersChange }")
  .ComponentSearch__tail
    dmc-button(label="submit" onPat="{ handleSubmitButtonPat }")

  script.
    import '../../organisms/dmc-parameters/index.tag';
    import script from './search';
    this.external(script);
