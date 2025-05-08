package thunder

import (
    "github.com/octoberswimmer/masc"
    "github.com/octoberswimmer/masc/elem"
    "github.com/octoberswimmer/masc/event"
)

// ButtonVariant defines the SLDS button style variant.
type ButtonVariant string

const (
    // VariantNeutral is the default neutral button style.
    VariantNeutral ButtonVariant = "slds-button_neutral"
    // VariantBrand is the brand button style.
    VariantBrand ButtonVariant = "slds-button_brand"
    // VariantDestructive is the destructive button style.
    VariantDestructive ButtonVariant = "slds-button_destructive"
)

// Button renders an SLDS button with the given label, style variant, and click handler.
// If variant is empty, VariantNeutral will be used.
func Button(label string, variant ButtonVariant, onClick func(*masc.Event)) masc.ComponentOrHTML {
    v := string(variant)
    if v == "" {
        v = string(VariantNeutral)
    }
    return elem.Button(
        masc.Markup(
            masc.Class("slds-button", v),
            event.Click(onClick),
        ),
        masc.Text(label),
    )
}